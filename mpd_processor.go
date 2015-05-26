package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	/**
	 * Any gap/overlap within a SegmentTimeline that is greater than or equal to
	 * this value (in seconds) will generate a warning message.
	 * @const {number}
	 */
	GAP_OVERLAP_WARNING_THRESHOLD = 1 / 32

	/**
	 * The maximum span, in seconds, that a SegmentIndex must account for when that
	 * SegmentIndex is being generated via a segment duration.
	 * @const {number}
	 */
	MAX_SEGMENT_INDEX_SPAN = 2 * 60

	/**
	 * The default value, in seconds, for MPD@minBufferTime if this attribute is
	 * missing.
	 * @const {number}
	 */
	DEFAULT_MIN_BUFFER_TIME = 5
)

/**
 * Creates an MpdProcessor, which validates MPDs, calculates start/duration
 * attributes, removes invalid Representations, and ultimately generates a
 * ManifestInfo.
 *
 * @param {ContentProtectionCallback} interpretContentProtection
 *
 */

type MpdProcessor struct {
	ManifestInfo ManifestInfo
}

func NewMpdProcessor() MpdProcessor {
	return MpdProcessor{}
}

/**
 * Processes the given MPD. Sets |this.periodInfos|.
 *
 * @param {Mpd} mpd
 */
func (mpdProcessor *MpdProcessor) Process(mpd *Mpd) {
	mpdProcessor.ManifestInfo = NewManifestInfo()
	mpdProcessor.validateSegmentInfo(mpd)
	mpdProcessor.calculateDurations(mpd)
	mpdProcessor.filterPeriods(mpd)
	mpdProcessor.createManifestInfo(*mpd)
}

/**
 * Ensures that each Representation has either a SegmentBase, SegmentList, or
 * SegmentTemplate.
 *
 * @param {Mpd} mpd
 */
func (mpdProcessor *MpdProcessor) validateSegmentInfo(mpd *Mpd) {
	for _, period := range mpd.Periods {
		for _, adaptationSet := range period.AdaptationSets {
			if adaptationSet.ContentType == "text" {
				continue
			}

			for k := 0; k < len(adaptationSet.Representations); k++ {
				representation := adaptationSet.Representations[k]

				n := 0
				if representation.SegmentBase != nil {
					n += 1
				}
				if representation.SegmentList != nil {
					n += 1
				}
				if representation.SegmentTemplate != nil {
					n += 1
				}

				if n == 0 {
					fmt.Printf("Representation does not contain any segment information.\r\nA Representation must contain one of SegmentBase, SegmentList, or SegmentTemplate. %s\r\n", representation)
					adaptationSet.Representations = append(adaptationSet.Representations[:k], adaptationSet.Representations[k+1:]...)
					k--
				} else if n != 1 {
					fmt.Printf("Representation contains multiple segment information sources.\r\nA Representation should only contain one of SegmentBase, SegmenstList, or SegmentTemplate. %s\r\n", representation)
					if representation.SegmentBase != nil {
						fmt.Printf("Using SegmentBase by default.")
						representation.SegmentList = nil
						representation.SegmentTemplate = nil
					} else if representation.SegmentList != nil {
						fmt.Printf("Using SegmentList by default.")
						representation.SegmentTemplate = nil
					} else {
						// asserts.unreachable();
						fmt.Println("unreachable")
					}
				}
			} // for k
		}
	}
}

/**
 * Attempts to calculate each Period's start attribute and duration attribute,
 * and attempts to calcuate the MPD's mediaPresentationDuration attribute.
 *
 * @see ISO/IEC 23009-1:2014 section 5.3.2.1
 *
 * @param {Mpd} mpd
 */
func (mpdProcessor *MpdProcessor) calculateDurations(mpd *Mpd) {
	if len(mpd.Periods) == 0 {
		return
	}

	if mpd.Periods[0].Start == -1 {
		mpd.Periods[0].Start = 0
	}

	// @mediaPresentationDuration should only be used if the MPD is static.
	if mpd.Type != "static" {
		mpd.MediaPresentationDuration = -1
	}

	// If there is only one Period then infer its duration.
	if mpd.MediaPresentationDuration != -1 && len(mpd.Periods) == 1 && mpd.Periods[0].Duration == -1 {
		mpd.Periods[0].Duration = mpd.MediaPresentationDuration
	}

	totalDuration := 0

	// True if |totalDuration| includes all periods, false if it only includes up
	// to the last Period in which a start time and duration could be
	// ascertained.
	totalDurationIncludesAllPeriods := true

	for i := 0; i < len(mpd.Periods); i++ {
		var previousPeriod *Period = nil
		var nextPeriod *Period = nil

		if i > 0 {
			previousPeriod = mpd.Periods[i-1]
		}

		period := mpd.Periods[i]

		// "The Period extends until the Period.start of the next Period, or until
		// the end of the Media Presentation in the case of the last Period."

		// nextPeriod = mpd.Periods[i + 1] || { start: mpd.MediaPresentationDuration }

		if (i + 1) < len(mpd.Periods) {
			nextPeriod = mpd.Periods[i+1]
		} else {
			nextPeriod = NewPeriod().(*Period)
			nextPeriod.Start = mpd.MediaPresentationDuration
		}

		// "If the 'start' attribute is absent, but the previous period contains a
		// 'duration' attribute, the start time of the new Period is the sum of the
		// start time of the previous period Period.start and the value of the
		// attribute 'duration' of the previous Period."
		if period.Start == -1 && previousPeriod != nil && previousPeriod.Start != -1 && previousPeriod.Duration != -1 {
			period.Start = previousPeriod.Start + previousPeriod.Duration
		}

		// "The difference between the start time of a Period and the start time
		// of the following Period is the duration of the media content represented
		// by this Period."
		if period.Duration == -1 && nextPeriod.Start != -1 {
			period.Duration = nextPeriod.Start - period.Start
		}

		if period.Start != -1 && period.Duration != -1 {
			totalDuration += period.Duration
		} else {
			totalDurationIncludesAllPeriods = false
		}
	}

	// "The Media Presentation Duration is provided either as the value of MPD
	// 'mediaPresentationDuration' attribute if present, or as the sum of
	// Period.start + Period.duration of the last Period."
	if mpd.MediaPresentationDuration != -1 {
		if mpd.MediaPresentationDuration != totalDuration {
			fmt.Println("@mediaPresentationDuration does not match the total duration of all periods.")
			// Assume mpd.mediaPresentationDuration is correct;
			// |totalDurationIncludesAllPeriods| may be false.
		}
	} else {
		finalPeriod := mpd.Periods[len(mpd.Periods)-1]
		if totalDurationIncludesAllPeriods {
			assert(finalPeriod.Start == -1 && finalPeriod.Duration == -1)
			assert(totalDuration == finalPeriod.Start+finalPeriod.Duration)
			mpd.MediaPresentationDuration = totalDuration
		} else {
			if finalPeriod.Start != -1 && finalPeriod.Duration != -1 {
				fmt.Println("Some Periods may not have valid start times or durations.")
				mpd.MediaPresentationDuration = finalPeriod.Start + finalPeriod.Duration
			} else {
				// Fallback to what we were able to compute.
				if mpd.Type == "static" {
					fmt.Println("Some Periods may not have valid start times or durations; @mediaPresentationDuration may not include the duration of all periods.")
					mpd.MediaPresentationDuration = totalDuration
				}
			}
		}
	}
}

/**
 * Removes invalid Representations from |mpd|.
 */
func (mpdProcessor *MpdProcessor) filterPeriods(mpd *Mpd) {
	for _, period := range mpd.Periods {
		for j := 0; j < len(period.AdaptationSets); j++ {
			adaptationSet := period.AdaptationSets[j]
			mpdProcessor.filterAdaptationSet(adaptationSet)
			if len(adaptationSet.Representations) == 0 {
				// Drop any AdaptationSet that is empty.
				// An error has already been logged.
				period.AdaptationSets = append(period.AdaptationSets[:j], period.AdaptationSets[j+1:]...)
				j--
			}
		}
	}
}

/**
 * Removes any Representation from the given AdaptationSet that has a different
 * MIME type than the MIME type of the first Representation of the
 * AdaptationSet.
 *
 * @param {AdaptationSet} adaptationSet
 */
func (mpdProcessor *MpdProcessor) filterAdaptationSet(adaptationSet *AdaptationSet) {
	desiredMimeType := ""

	for i := 0; i < len(adaptationSet.Representations); i++ {
		representation := adaptationSet.Representations[i]
		mimeType := representation.MimeType

		if desiredMimeType == "" {
			desiredMimeType = mimeType
		} else if mimeType != desiredMimeType {
			fmt.Printf("Representation has an inconsistent mime type. %s\r\n", adaptationSet.Representations[i])
			adaptationSet.Representations = append(adaptationSet.Representations[:i], adaptationSet.Representations[i+1:]...)
			i--
		}
	}
}

/**
 * Creates a ManifestInfo from |mpd|.
 *
 * @param {Mpd} mpd
 */
func (mpdProcessor *MpdProcessor) createManifestInfo(mpd Mpd) {
	mpdProcessor.ManifestInfo.Live = (mpd.Type == "dynamic")
	mpdProcessor.ManifestInfo.MinBufferTime = mpd.MinBufferTime

	for i := 0; i < len(mpd.Periods); i++ {
		period := mpd.Periods[i]

		periodInfo := NewPeriodInfo()
		periodInfo.Id = period.Id

		if period.Start == -1 {
			periodInfo.Start = 0
		} else {
			periodInfo.Start = period.Start
		}

		periodInfo.Duration = period.Duration

		for _, adaptationSet := range period.AdaptationSets {

			streamSetInfo := NewStreamSetInfo()
			streamSetInfo.Id = adaptationSet.Id
			streamSetInfo.Main = adaptationSet.Main
			streamSetInfo.ContentType = adaptationSet.ContentType
			streamSetInfo.Lang = adaptationSet.Lang

			// Keep track of the largest end time of all segment references so that
			// we can set a Period duration if one was not explicitly set in the MPD
			// or calculated from calculateDurations_().
			maxLastEndTime := 0

			for _, representation := range adaptationSet.Representations {

				// Get common DRM schemes.
				// var commonDrmSchemes = streamSetInfo.drmSchemes.slice(0);
				// this.updateCommonDrmSchemes_(representation, commonDrmSchemes);
				// if (commonDrmSchemes.length == 0 &&
				//     streamSetInfo.drmSchemes.length > 0) {
				//   shaka.log.warning(
				//       'Representation does not contain any DRM schemes that are in ' +
				//       'common with other Representations within its AdaptationSet.',
				//       representation);
				//   continue;
				// }

				streamInfo := mpdProcessor.createStreamInfo(mpd, *period, *representation)
				if streamInfo == nil {
					// An error has already been logged.
					continue
				}

				streamSetInfo.StreamInfos = append(streamSetInfo.StreamInfos, streamInfo)
				// streamSetInfo.drmSchemes = commonDrmSchemes;

				if streamInfo.SegmentIndex != nil && streamInfo.SegmentIndex.Length() > 0 {
					maxLastEndTime = Max(maxLastEndTime, streamInfo.SegmentIndex.Last().EndTime)
				}
			}

			periodInfo.StreamSetInfos = append(periodInfo.StreamSetInfos, streamSetInfo)

			if periodInfo.Duration == -1 {
				periodInfo.Duration = maxLastEndTime

				// If the MPD is dynamic then the Period's duration will likely change
				// after we re-process/update the MPD. When the Period's duration
				// changes we must update the MediaSource object that is presenting the
				// MPD's content, so we can append new media segments to the
				// MediaSource's SourceBuffers. However, changing the MediaSource's
				// duration is challenging as it requires synchronizing the states of
				// multiple SourceBuffers. We can leave the Period's duration as
				// undefined, but then we cannot seek (even programmatically).
				//
				// So, if the MPD is dynamic just set the Period's duration to a
				// "large" value. This ensures that we can seek and that we can append
				// new media segments. This does cause a poor UX if we use the video
				// element's default controls, but we shouldn't use the default
				// controls for live anyways.
				//
				// TODO: Remove this hack once SourceBuffer synchronization is
				// implemented.
				if mpd.Type == "dynamic" {
					periodInfo.Duration += 60 * 60 * 24 * 30
				}
			}
		}

		mpdProcessor.ManifestInfo.PeriodInfos = append(mpdProcessor.ManifestInfo.PeriodInfos, periodInfo)
	}
}

// /**
//  * Updates |commonDrmSchemes|.
//  *
//  * If |commonDrmSchemes| is empty then after this function is called
//  * |commonDrmSchemes| will equal |representation|'s application provided DRM
//  * schemes.
//  *
//  * Otherwise, if |commonDrmSchemes| is non-empty then after this function is
//  * called |commonDrmSchemes| will equal the intersection between
//  * |representation|'s application provided DRM schemes and |commonDrmSchemes|
//  * at the time this function was called.
//  *
//  * @param {!shaka.dash.mpd.Representation} representation
//  * @param {!Array.<!shaka.player.DrmSchemeInfo>} commonDrmSchemes
//  *
//  * @private
//  */
// // shaka.dash.MpdProcessor.prototype.updateCommonDrmSchemes_ = function(
// //     representation, commonDrmSchemes) {
// //   var drmSchemes = this.getDrmSchemeInfos_(representation);

// //   if (commonDrmSchemes.length == 0) {
// //     Array.prototype.push.apply(commonDrmSchemes, drmSchemes);
// //     return;
// //   }

// //   for (var i = 0; i < commonDrmSchemes.length; ++i) {
// //     var found = false;
// //     for (var j = 0; j < drmSchemes.length; ++j) {
// //       if (commonDrmSchemes[i].key() == drmSchemes[j].key()) {
// //         found = true;
// //         break;
// //       }
// //     }
// //     if (!found) {
// //       commonDrmSchemes.splice(i, 1);
// //       --i;
// //     }
// //   }
// // };

// /**
//  * Gets the application provided DrmSchemeInfos for the given Representation.
//  *
//  * @param {!shaka.dash.mpd.Representation} representation
//  * @return {!Array.<!shaka.player.DrmSchemeInfo>} The application provided
//  *     DrmSchemeInfos. A dummy scheme, which has an empty |keySystem| string,
//  *     is used for unencrypted content.
//  * @private
//  */
// // shaka.dash.MpdProcessor.prototype.getDrmSchemeInfos_ =
// //     function(representation) {
// //   var drmSchemes = [];
// //   if (representation.contentProtections.length == 0) {
// //     // Return a single item which indicates that the content is unencrypted.
// //     drmSchemes.push(shaka.player.DrmSchemeInfo.createUnencrypted());
// //   } else if (this.interpretContentProtection_) {
// //     for (var i = 0; i < representation.contentProtections.length; ++i) {
// //       var contentProtection = representation.contentProtections[i];
// //       var drmSchemeInfo = this.interpretContentProtection_(contentProtection);
// //       if (drmSchemeInfo) {
// //         drmSchemes.push(drmSchemeInfo);
// //       }
// //     }
// //   }
// //   return drmSchemes;
// // };

/**
 * Creates a StreamInfo from the given Representation.
 *
 * @param {Mpd} mpd
 * @param {Period} period
 * @param {Representation} representation
 * @return {StreamInfo} The new StreamInfo on success; otherwise,
 *     return null.
 */
func (mpdProcessor *MpdProcessor) createStreamInfo(mpd Mpd, period Period, representation Representation) *StreamInfo {
	streamInfo := NewStreamInfo()

	streamInfo.Id = representation.Id
	streamInfo.Bandwidth = representation.Bandwidth
	streamInfo.Width = representation.Width
	streamInfo.Height = representation.Height
	streamInfo.MimeType = representation.MimeType
	streamInfo.Codecs = representation.Codecs

	ok := false

	if representation.SegmentBase != nil {
		fmt.Println("BuildStreamInfoFromSegmentBase")
		ok = mpdProcessor.buildStreamInfoFromSegmentBase(representation.SegmentBase, &streamInfo)
	} else if representation.SegmentList != nil {
		fmt.Println("BuildStreamInfoFromSegmentList")
		ok = mpdProcessor.buildStreamInfoFromSegmentList(representation.SegmentList, &streamInfo)
	} else if representation.SegmentTemplate != nil {
		ok = mpdProcessor.buildStreamInfoFromSegmentTemplate(mpd, period, representation, &streamInfo)
	} else if strings.Split(representation.MimeType, "/")[0] == "text" {
		// All we need is a URL for subtitles.
		streamInfo.MediaUrl = representation.BaseUrl.Url
		ok = true
	} else {
		fmt.Println("unreachable")
		// shaka.asserts.unreachable();
	}

	if ok {
		return &streamInfo
	} else {
		return nil
	}
}

/**
 * Builds a StreamInfo from a SegmentBase.
 *
 * @param {SegmentBase} segmentBase
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromSegmentBase(segmentBase *SegmentBase, streamInfo *StreamInfo) bool {

	assert(segmentBase.Timescale > 0)

	hasSegmentIndexMetadata := segmentBase.IndexRange != nil || (segmentBase.RepresentationIndex != nil && segmentBase.RepresentationIndex.Range != nil)
	if !hasSegmentIndexMetadata || segmentBase.BaseUrl == nil {
		fmt.Printf("A SegmentBase must have a segment index URL and a base URL. %s\r\n", segmentBase)
		return false
	}

	if segmentBase.PresentationTimeOffset != -1 {
		// Each timestamp within each media segment is relative to the start of the
		// Period minus @presentationTimeOffset. So to align the start of the first
		// segment to the start of the Period we must apply an offset of -1 *
		// @presentationTimeOffset seconds to each timestamp within each media
		// segment.
		streamInfo.TimestampOffset = -1 * segmentBase.PresentationTimeOffset / segmentBase.Timescale
	}

	// If a RepresentationIndex does not exist then fallback to the indexRange
	// attribute.
	representationIndex := segmentBase.RepresentationIndex
	if representationIndex == nil {
		representationIndex = &RepresentationIndex{}
		representationIndex.Url = segmentBase.BaseUrl.Url
		if segmentBase.IndexRange != nil {
			representationIndex.Range = segmentBase.IndexRange.Clone()
		} else {
			representationIndex.Range = nil
		}
	}

	// Set StreamInfo properties.
	streamInfo.MediaUrl = segmentBase.BaseUrl.Url

	segmentIndexInfo, _ := mpdProcessor.createSegmentMetadataInfo(representationIndex)
	streamInfo.SegmentIndexInfo = &segmentIndexInfo

	segmentInitializationInfo, _ := mpdProcessor.createSegmentMetadataInfo(segmentBase.Initialization)
	streamInfo.SegmentInitializationInfo = &segmentInitializationInfo

	return true
}

/**
 * Builds a StreamInfo from a SegmentList.
 *
 * @param {SegmentList} segmentList
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 * @private
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromSegmentList(segmentList *SegmentList, streamInfo *StreamInfo) bool {
	assert(segmentList.Timescale > 0)

	if segmentList.SegmentDuration == -1 && len(segmentList.SegmentUrls) > 1 {
		fmt.Printf("A SegmentList without a segment duration can only have one segment.%s\r\n", segmentList)
		return false
	}

	segmentInitializationInfo, _ := mpdProcessor.createSegmentMetadataInfo(segmentList.Initialization)
	streamInfo.SegmentInitializationInfo = &segmentInitializationInfo

	lastEndTime := 0
	references := make([]*SegmentReference, 0)

	for i, segmentUrl := range segmentList.SegmentUrls {

		// Compute the segment's unscaled start time.
		var startTime int
		if i == 0 {
			startTime = 0
		} else {
			startTime = lastEndTime
		}
		assert(startTime >= 0)

		endTime := -1
		scaledEndTime := -1

		scaledStartTime := startTime / segmentList.Timescale

		// If segmentList.segmentDuration is null then there must only be one
		// segment.
		if segmentList.SegmentDuration != -1 {
			endTime = startTime + segmentList.SegmentDuration
			scaledEndTime = endTime / segmentList.Timescale
		}

		lastEndTime = endTime

		startByte := 0
		endByte := 0

		if segmentUrl.MediaRange != nil {
			startByte = segmentUrl.MediaRange.Begin
			endByte = segmentUrl.MediaRange.End
		}

		segmentReference := NewSegmentReference(
			startTime,
			scaledStartTime,
			scaledEndTime,
			startByte,
			endByte,
			segmentUrl.MediaUrl)

		references = append(references, &segmentReference)

	}

	// Set StreamInfo properties.
	segmentIndex := NewSegmentIndex(references)
	streamInfo.SegmentIndex = &segmentIndex
	fmt.Printf("Generated SegmentIndex from SegmentList %s\r\n", streamInfo.SegmentIndex)

	return true
}

/**
 * Builds a StreamInfo from a SegmentTemplate.
 *
 * @param {Mpd} mpd
 * @param {Period} period
 * @param {Representation} representation
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromSegmentTemplate(mpd Mpd, period Period, representation Representation, streamInfo *StreamInfo) bool {

	assert(representation.SegmentTemplate != nil)

	segmentTemplate := representation.SegmentTemplate

	ok := false

	// Prefer an explicit segment index URL, then a SegmentTimeline, and then a
	// segment duration.
	if segmentTemplate.IndexUrlTemplate != "" {
		if segmentTemplate.Timeline != nil {
			fmt.Printf("Ignoring SegmentTimeline because an explicit segment index URL was provided for the SegmentTemplate.%s\r\n", representation)
		}
		if segmentTemplate.SegmentDuration != -1 {
			fmt.Printf("Ignoring segment duration because an explicit segment index URL was provided for the SegmentTemplate.%s\r\n", representation)
		}
		fmt.Println("BuildStreamInfoFromIndexUrlTemplate")
		ok = mpdProcessor.buildStreamInfoFromIndexUrlTemplate(representation, streamInfo)
	} else if segmentTemplate.Timeline != nil {
		if segmentTemplate.SegmentDuration != -1 {
			fmt.Printf("Ignoring segment duration because a SegmentTimeline was provided for the SegmentTemplate.%s\r\n", representation)
		}
		fmt.Println("BuildStreamInfoFromSegmentTimeline")
		ok = mpdProcessor.buildStreamInfoFromSegmentTimeline(mpd, period, representation, streamInfo)
	} else if segmentTemplate.SegmentDuration != -1 {
		fmt.Println("BuildStreamInfoFromSegmentDuration")
		ok = mpdProcessor.buildStreamInfoFromSegmentDuration(mpd, period, representation, streamInfo)
	} else {
		fmt.Printf("SegmentTemplate does not provide an explicit segment index URL, a SegmentTimeline, or a segment duration.%s\r\n", representation)
		ok = false
	}

	return ok
}

/**
 * Builds a StreamInfo from a SegmentTemplate with an index URL template.
 *
 * @param {Representation} representation
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 * @private
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromIndexUrlTemplate(representation Representation, streamInfo *StreamInfo) bool {
	fmt.Println("buildStreamInfoFromIndexUrlTemplate")

	assert(representation.SegmentTemplate != nil)
	assert(representation.SegmentTemplate.IndexUrlTemplate != "")
	assert(representation.SegmentTemplate.Timescale > 0)

	segmentTemplate := representation.SegmentTemplate

	// Generate the media URL. Since there is no SegmentTimeline there is only
	// one media URL, so just map $Number$ to 1 and $Time$ to 0.
	mediaUrl := ""

	if segmentTemplate.MediaUrlTemplate != "" {
		filledUrlTemplate := mpdProcessor.fillUrlTemplate(segmentTemplate.MediaUrlTemplate, representation.Id, 1, representation.Bandwidth, 0)

		if filledUrlTemplate == "" {
			// An error has already been logged.
			return false
		}

		if representation.BaseUrl != nil {
			mediaUrl = representation.BaseUrl.Url + filledUrlTemplate
		} else {
			mediaUrl = filledUrlTemplate
		}
		// mediaUrl = representation.baseUrl ?
		//            representation.baseUrl.resolve(filledUrlTemplate) :
		//            filledUrlTemplate;
	} else {
		// Fallback to the Representation's URL.
		mediaUrl = representation.BaseUrl.Url
	}

	// Generate a RepresentationIndex.
	var err error
	var representationIndex RepresentationIndex

	if representationIndex, err = mpdProcessor.generateRepresentationIndex(representation); err != nil {
		// An error has already been logged.
		return false
	}

	// Generate an Initialization.
	var initialization Initialization
	if segmentTemplate.InitializationUrlTemplate != "" {
		if initialization, err = mpdProcessor.generateInitialization(representation); err != nil {
			// An error has already been logged.
			return false
		}
	}

	// Set StreamInfo properties.
	streamInfo.MediaUrl = mediaUrl

	if segmentTemplate.PresentationTimeOffset != -1 {
		streamInfo.TimestampOffset = -1 * segmentTemplate.PresentationTimeOffset / segmentTemplate.Timescale
	}

	if segmentIndexInfo, err := mpdProcessor.createSegmentMetadataInfo(&representationIndex); err != nil {
		streamInfo.SegmentIndexInfo = nil
	} else {
		streamInfo.SegmentIndexInfo = &segmentIndexInfo
	}

	if segmentInitializationInfo, err := mpdProcessor.createSegmentMetadataInfo(&initialization); err != nil {
		streamInfo.SegmentInitializationInfo = nil
	} else {
		streamInfo.SegmentInitializationInfo = &segmentInitializationInfo
	}

	return true
}

/**
 * Generates a RepresentationIndex from a SegmentTemplate.
 *
 * @param {Representation} representation
 * @return {RepresentationIndex} A RepresentationIndex on
 *     success, null if no index URL template exists or an error occurred.
 */
func (mpdProcessor *MpdProcessor) generateRepresentationIndex(representation Representation) (RepresentationIndex, error) {
	assert(representation.SegmentTemplate != nil)

	representationIndex := RepresentationIndex{}

	segmentTemplate := representation.SegmentTemplate
	assert(segmentTemplate.IndexUrlTemplate != "")
	if segmentTemplate.IndexUrlTemplate == "" {
		return representationIndex, errors.New("missing index url template")
	}

	// $Number$ and $Time$ cannot be present in an index URL template.
	filledUrlTemplate := mpdProcessor.fillUrlTemplate(segmentTemplate.IndexUrlTemplate, representation.Id, 0, representation.Bandwidth, 0)

	if filledUrlTemplate == "" {
		// An error has already been logged.
		return representationIndex, errors.New("missing filled url template")
	}

	if representation.BaseUrl != nil && filledUrlTemplate != "" {
		representationIndex.Url = representation.BaseUrl.Url + filledUrlTemplate
	} else {
		representationIndex.Url = filledUrlTemplate
	}

	return representationIndex, nil
}

/**
 * Builds a StreamInfo from a SegmentTemplate with a SegmentTimeline.
 *
 * @param {Mpd} mpd
 * @param {Period} period
 * @param {Representation} representation
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromSegmentTimeline(mpd Mpd, period Period, representation Representation, streamInfo *StreamInfo) bool {
	assert(representation.SegmentTemplate != nil)
	assert(representation.SegmentTemplate.Timeline != nil)
	assert(representation.SegmentTemplate.Timescale > 0)

	if period.Start == -1 {
		fmt.Printf("Cannot instantiate SegmentTemplate: the period's start time is unknown.%s\r\n", representation)
		return false
	}

	segmentTemplate := representation.SegmentTemplate
	if segmentTemplate.MediaUrlTemplate == "" {
		fmt.Printf("Cannot instantiate SegmentTemplate: SegmentTemplate does not have a media URL template.%s\r\n", representation)
		return false
	}

	timeline := mpdProcessor.createTimeline(segmentTemplate)
	if timeline == nil {
		// An error has already been logged.
		return false
	}

	// Compute the earliest available timestamp. Assume the MPD only contains
	// segments that are available. This simplifies the calculation below by
	// allowing us to ignore @availabilityStartTime. If we did use
	// @availabilityStartTime then the calculation below would be more
	// complicated than the calculations in computeAvailableSegmentRange_() since
	// the duration of each segment is variable here.
	earliestAvailableTimestamp := 0
	if mpd.Type == "dynamic" && len(timeline) > 0 {
		index := Max(0, len(timeline)-2)
		timeShiftBufferDepth := mpd.TimeShiftBufferDepth
		earliestAvailableTimestamp = (timeline[index].Start / segmentTemplate.Timescale) - timeShiftBufferDepth
	}

	// Generate a SegmentIndex.
	references := make([]*SegmentReference, 0)

	for i := 0; i < len(timeline); i++ {
		startTime := timeline[i].Start
		endTime := timeline[i].End

		// Compute the segment's scaled start time and scaled end time.
		scaledStartTime := startTime / segmentTemplate.Timescale
		scaledEndTime := endTime / segmentTemplate.Timescale

		if scaledStartTime < earliestAvailableTimestamp {
			// Skip unavailable segments.
			continue
		}

		absoluteSegmentNumber := i + segmentTemplate.StartNumber

		// Compute the media URL template placeholder replacements.
		segmentReplacement := absoluteSegmentNumber
		timeReplacement := startTime

		// Generate the media URL.
		assert(segmentTemplate.MediaUrlTemplate != "")
		filledUrlTemplate := mpdProcessor.fillUrlTemplate(
			segmentTemplate.MediaUrlTemplate,
			representation.Id,
			segmentReplacement,
			representation.Bandwidth,
			timeReplacement)

		if filledUrlTemplate == "" {
			// An error has already been logged.
			return false
		}

		mediaUrl := ""
		if representation.BaseUrl != nil {
			mediaUrl = representation.BaseUrl.Url + filledUrlTemplate
		} else {
			mediaUrl = filledUrlTemplate
		}
		segmentReference := NewSegmentReference(startTime, scaledStartTime, scaledEndTime, 0 /* startByte */, -1 /* endByte */, mediaUrl)
		references = append(references, &segmentReference)
	}

	// Generate an Initialization. If there are no references then assume that
	// the intialization segment is not available.
	var initialization Initialization
	var err error
	if segmentTemplate.InitializationUrlTemplate != "" && len(references) > 0 {
		if initialization, err = mpdProcessor.generateInitialization(representation); err != nil {
			// An error has already been logged.
			return false
		}
	}

	// Set StreamInfo properties.
	if segmentTemplate.PresentationTimeOffset != -1 {
		streamInfo.TimestampOffset = -1 * segmentTemplate.PresentationTimeOffset / segmentTemplate.Timescale
	}

	if mpd.Type == "dynamic" && len(references) > 0 {
		minBufferTime := mpdProcessor.ManifestInfo.MinBufferTime
		bestAvailableTimestamp := references[len(references)-1].StartTime - minBufferTime

		if bestAvailableTimestamp >= earliestAvailableTimestamp {
			fmt.Printf("The best available segment is still available!")
		} else {
			// NOTE: @minBufferTime is large compared to @timeShiftBufferDepth, so we
			// can't start as far back, for buffering, as we'd like.
			bestAvailableTimestamp = earliestAvailableTimestamp
			fmt.Printf("The best available segment is no longer available.")
		}

		for i := 0; i < len(references); i++ {
			if references[i].EndTime >= bestAvailableTimestamp {
				streamInfo.CurrentSegmentStartTime = references[i].StartTime
				break
			}
		}

		assert(streamInfo.CurrentSegmentStartTime != -1)
	}

	if segmentInitializationInfo, err := mpdProcessor.createSegmentMetadataInfo(&initialization); err != nil {
		streamInfo.SegmentInitializationInfo = nil
	} else {
		streamInfo.SegmentInitializationInfo = &segmentInitializationInfo
	}

	segmentIndex := NewSegmentIndex(references)
	streamInfo.SegmentIndex = &segmentIndex
	fmt.Printf("Generated SegmentIndex from SegmentTimeline", streamInfo.SegmentIndex)

	return true
}

/**
 * Expands a SegmentTimeline into a simple array-based timeline.
 *
 * @return {Array.<{start: number, end: number}>}
 */
func (mpdProcessor *MpdProcessor) createTimeline(segmentTemplate *SegmentTemplate) []TimeLine {
	assert(segmentTemplate.Timeline != nil)

	timePoints := segmentTemplate.Timeline.TimePoints
	lastEndTime := 0

	/** @type {!Array.<{start: number, end: number}>} */
	timeline := make([]TimeLine, 0)

	for i := 0; i < len(timePoints); i++ {
		repeat := 0

		if timePoints[i].Repeat != -1 {
			repeat = timePoints[i].Repeat
		}

		for j := 0; j <= repeat; j++ {
			if timePoints[i].Duration == -1 {
				fmt.Printf("SegmentTimeline 'S' element does not have a duration.%s\r\n", timePoints[i])
				return nil
			}

			// Compute the segment's unscaled start time and unscaled end time.
			var startTime int
			if timePoints[i].StartTime != -1 && j == 0 {
				startTime = timePoints[i].StartTime
			} else {
				if i == 0 && j == 0 {
					startTime = 0
				} else {
					startTime = lastEndTime
				}
			}
			assert(startTime >= 0)
			endTime := startTime + timePoints[i].Duration

			// The end of the last segment may end before the start of the current
			// segment (a gap) or may end after the start of the current segment (an
			// overlap). If there is a gap/overlap then stretch/compress the end of
			// the last segment to the start of the current segment.
			//
			// Note: it is possible to move the start of the current segment to the
			// end of the last segment, but this complicates the computation of the
			// $Time$ placeholder.
			if len(timeline) > 0 && startTime != lastEndTime {
				delta := startTime - lastEndTime

				if Abs(delta/segmentTemplate.Timescale) >= GAP_OVERLAP_WARNING_THRESHOLD {
					fmt.Printf("SegmentTimeline contains a large gap/overlap, the content may have errors in it.%s\r\n", timePoints[i])
				}

				timeline[len(timeline)-1].End = startTime
			}

			lastEndTime = endTime

			timeline = append(timeline, NewTimeLine(startTime, endTime))
		} // for j
	}

	return timeline
}

/**
 * Builds a StreamInfo from a SegmentTemplate with a segment duration.
 *
 * @param {Mpd} mpd
 * @param {Period} period
 * @param {Representation} representation
 * @param {StreamInfo} streamInfo
 * @return {boolean} True on success.
 */
func (mpdProcessor *MpdProcessor) buildStreamInfoFromSegmentDuration(mpd Mpd, period Period, representation Representation, streamInfo *StreamInfo) bool {
	assert(representation.SegmentTemplate != nil)
	assert(representation.SegmentTemplate.SegmentDuration != -1)
	assert(representation.SegmentTemplate.Timescale > 0)

	if period.Start == -1 {
		fmt.Printf("Cannot instantiate SegmentTemplate: the period's start time is unknown.%s\r\n", representation)
		return false
	}

	segmentTemplate := representation.SegmentTemplate
	if segmentTemplate.MediaUrlTemplate == "" {
		fmt.Printf("Cannot instantiate SegmentTemplate: SegmentTemplate does not have a media URL template.%s\r\n", representation)
		return false
	}

	// The number of segment references to generate starting from the earliest
	// available segment to the current segment, but not counting the current
	// segment.
	numSegmentsBeforeCurrentSegment := 0

	// Find the earliest available segment and the current segment. All segment
	// numbers are relative to the start of |period| unless marked otherwise.
	var earliestSegmentNumber int = -1
	var currentSegmentNumber int = -1

	// TODO: handle dynamic
	// if mpd.Type == "dynamic" {
	// 	pair := mpdProcessor.ComputeAvailableSegmentRange(mpd, period, *segmentTemplate)
	// 	if pair != nil {
	// 		// Build the SegmentIndex starting from the earliest available segment.
	// 		earliestSegmentNumber = pair.Earliest
	// 		currentSegmentNumber = pair.Current
	// 		numSegmentsBeforeCurrentSegment = currentSegmentNumber - earliestSegmentNumber
	// 		assert(numSegmentsBeforeCurrentSegment >= 0)
	// 	}
	// } else {
	// 	earliestSegmentNumber = 1
	// }
	// END OF TODO
	earliestSegmentNumber = 1 // REMOVE THIS LINE WHEN HANDLING DYNAMIC

	assert(earliestSegmentNumber == -1 || earliestSegmentNumber >= 0)

	// The optimal number of segment references to generate starting from, and
	// including, the current segment
	numSegmentsFromCurrentSegment := 0

	// Note that if |earliestSegmentNumber| is undefined then the current segment
	// is not available.
	if earliestSegmentNumber >= 0 {
		numSegmentsFromCurrentSegment = mpdProcessor.computeOptimalSegmentIndexSize(mpd, period, *segmentTemplate)
		if numSegmentsFromCurrentSegment == -1 {
			// An error has already been logged.
			return false
		}
	}

	fmt.Printf("numSegmentsBeforeCurrentSegment: %d, numSegmentsFromCurrentSegment: %d\r\n", numSegmentsBeforeCurrentSegment, numSegmentsFromCurrentSegment)
	totalNumSegments := numSegmentsBeforeCurrentSegment + numSegmentsFromCurrentSegment
	references := make([]*SegmentReference, 0)

	for i := 0; i < totalNumSegments; i++ {
		segmentNumber := i + earliestSegmentNumber

		startTime := (segmentNumber - 1) * segmentTemplate.SegmentDuration
		endTime := startTime + segmentTemplate.SegmentDuration

		scaledStartTime := startTime / segmentTemplate.Timescale
		scaledEndTime := endTime / segmentTemplate.Timescale

		absoluteSegmentNumber := (segmentNumber - 1) + segmentTemplate.StartNumber

		// Compute the media URL template placeholder replacements.
		segmentReplacement := absoluteSegmentNumber
		timeReplacement := ((segmentNumber - 1) + (segmentTemplate.StartNumber - 1)) * segmentTemplate.SegmentDuration

		// Generate the media URL.
		assert(segmentTemplate.MediaUrlTemplate != "")
		var filledUrlTemplate = mpdProcessor.fillUrlTemplate(
			segmentTemplate.MediaUrlTemplate,
			representation.Id,
			segmentReplacement,
			representation.Bandwidth,
			timeReplacement)

		if filledUrlTemplate == "" {
			// An error has already been logged.
			return false
		}

		mediaUrl := ""
		if representation.BaseUrl != nil {
			mediaUrl = representation.BaseUrl.Url + filledUrlTemplate
		} else {
			mediaUrl = filledUrlTemplate
		}
		segmentRef := NewSegmentReference(startTime, scaledStartTime, scaledEndTime, 0 /* startByte */, -1 /* endByte */, mediaUrl)
		references = append(references, &segmentRef)
	}

	// Generate an Initialization. If there are no references then assume that
	// the intialization segment is not available.
	var initialization Initialization
	var err error

	if segmentTemplate.InitializationUrlTemplate != "" && len(references) > 0 {
		if initialization, err = mpdProcessor.generateInitialization(representation); err != nil {
			// An error has already been logged.
			return false
		}
	}

	// Set StreamInfo properties.
	if segmentTemplate.PresentationTimeOffset != -1 {
		streamInfo.TimestampOffset = -1 * segmentTemplate.PresentationTimeOffset / segmentTemplate.Timescale
	}

	if mpd.Type == "dynamic" && len(references) > 0 {
		assert(currentSegmentNumber == -1)
		scaledSegmentDuration := segmentTemplate.SegmentDuration / segmentTemplate.Timescale
		streamInfo.CurrentSegmentStartTime = (currentSegmentNumber - 1) * scaledSegmentDuration
	}

	if segmentMetadataInfo, err := mpdProcessor.createSegmentMetadataInfo((Node)(&initialization)); err != nil {
		streamInfo.SegmentInitializationInfo = nil
	} else {
		streamInfo.SegmentInitializationInfo = &segmentMetadataInfo
	}

	segmentIdx := NewSegmentIndex(references)
	streamInfo.SegmentIndex = &segmentIdx

	return true
}

/**
 * Computes the optimal number of segment references, N, for |period|.  If the
 * MPD is static then N * segmentDuration is the smallest multiple of
 * segmentDuration >= |period|'s duration; if the MPD is dynamic then N *
 * segmentDuration is the smallest multiple of segmentDuration >= the minimum
 * of |period|'s duration, minimumUpdatePeriod, and MAX_SEGMENT_INDEX_SPAN.
 *
 * If the MPD is dynamic, and at least one segment is available, then N can be
 * regarded as the number of segment references that we can generate right now,
 * such that the generated segment references will all be valid when it's time
 * to actually fetch the corresponding segments.
 *
 * @param {Mpd} mpd
 * @param {Period} period
 * @param {SegmentTemplate} segmentTemplate
 * @return {?number}
 * @private
 */
func (mpdProcessor *MpdProcessor) computeOptimalSegmentIndexSize(mpd Mpd, period Period, segmentTemplate SegmentTemplate) int {

	assert(segmentTemplate.SegmentDuration == -1)
	assert(segmentTemplate.Timescale > 0)

	var duration int = -1
	if mpd.Type == "static" {
		if period.Duration != -1 {
			duration = period.Duration
		} else {
			fmt.Printf(
				"Cannot instantiate SegmentTemplate: the Period's duration is unknown.%s\r\n", period)
			return -1
		}
	} else {
		// TODO: support dynamic!
		// Note that |period|'s duration and @minimumUpdatePeriod may be very
		// large, so fallback to a default value if necessary. The VideoSource is
		// responsible for generating new SegmentIndexes when it needs them.
		// END OF TODO, UNCOMMENT NEXT TWO LINES TO SUPPORT DYNAMIC
		// duration = Min(period.Duration || Number.POSITIVE_INFINITY, mpd.MinUpdatePeriod || Number.POSITIVE_INFINITY)
		// duration = Min(duration, MAX_SEGMENT_INDEX_SPAN)
	}
	// assert(duration && (duration != Number.POSITIVE_INFINITY), "duration should not be zero or infinity!")
	assert(duration == -1)

	scaledSegmentDuration := float64(segmentTemplate.SegmentDuration) / float64(segmentTemplate.Timescale)

	n := math.Ceil(float64(duration) / scaledSegmentDuration)
	assert(n >= 1)
	return int(n)
}

/**
 * Computes the segment numbers of the earliest segment and the current
 * segment, both relative to the start of |period|. Assumes the MPD is dynamic.
 * |segmentTemplate| must have a segment duration.
 *
 * The earliest segment is the segment with the smallest start time that is
 * still available from the media server. The current segment is the segment
 * with the largest start time that is available from the media server and that
 * also respects the suggestedPresentationDelay attribute and the minBufferTime
 * attribute.
 *
 * @param {!shaka.dash.mpd.Mpd} mpd
 * @param {!shaka.dash.mpd.Period} period
 * @param {!shaka.dash.mpd.SegmentTemplate} segmentTemplate
 * @return {?{earliest: number, current: number}} Two segment numbers, or null
 *     if the stream is not available yet.
 * @private
 */
// func (mpdProcessor *MpdProcessor) ComputeAvailableSegmentRange(mpd Mpd, period Period, segmentTemplate SegmentTemplate) *Pair {
// 	currentTime := shaka.util.Clock.now() / 1000.0
// 	var availabilityStartTime int

// 	if mpd.AvailabilityStartTime != -1 {
// 		availabilityStartTime = mpd.AvailabilityStartTime
// 	} else {
// 		availabilityStartTime = currentTime
// 	}

// 	if availabilityStartTime > currentTime {
// 		fmt.Printf("The stream is not available yet!%s\r\n", period)
// 		return nil
// 	}

// 	minBufferTime := mpd.minBufferTime
// 	suggestedPresentationDelay := mpd.SuggestedPresentationDelay

// 	// The following diagram shows the relationship between the values we use to
// 	// compute the current segment number; descriptions of each value are given
// 	// within the code. The diagram depicts the media presentation timeline. 0
// 	// corresponds to availabilityStartTime + period.start in wall-clock time,
// 	// and currentPresentationTime corresponds to currentTime in wall-clock time.
// 	//
// 	// Legend:
// 	// CPT: currentPresentationTime
// 	// EAT: earliestAvailableSegmentStartTime
// 	// LAT: latestAvailableSegmentStartTime
// 	// BAT: bestAvailableSegmentStartTime
// 	// SD:  scaledSegmentDuration.
// 	// SPD: suggestedPresentationDelay
// 	// MBT: minBufferTime
// 	// TSB: timeShiftBufferDepth
// 	//
// 	// Time:
// 	//   <---|-----------------+--------+-----------------+----------|--------->
// 	//       0                EAT      BAT               LAT        CPT
// 	//                                                      |---SD---|
// 	//                                      |-MBT-|--SPD--|
// 	//                      |---SD---|---SD---|<--------TSB--------->|
// 	// Segments:
// 	//   <---1--------2--------3--------4--------5--------6--------7--------8-->
// 	//       |---SD---|---SD---| ...

// 	assert(segmentTemplate.SegmentDuration != -1)
// 	assert(segmentTemplate.Timescale > 0)
// 	scaledSegmentDuration := segmentTemplate.SegmentDuration / segmentTemplate.Timescale

// 	// The current presentation time, which is the amount of time since the start
// 	// of the Period.
// 	currentPresentationTime := currentTime - (availabilityStartTime + period.Start)
// 	if currentPresentationTime < 0 {
// 		fmt.Printf("The Period is not available yet!%s\r\n", period)
// 		return nil
// 	}

// 	// Compute the segment start time of the earliest available segment, i.e.,
// 	// the segment that starts furthest from the present but is still available).
// 	// The MPD spec. indicates that
// 	//
// 	// SegmentAvailabilityStartTime =
// 	//   MpdAvailabilityStartTime + PeriodStart + SegmentStart + SegmentDuration
// 	//
// 	// SegmentAvailabilityEndTime =
// 	//   SegmentAvailabilityStartTime + SegmentDuration + TimeShiftBufferDepth
// 	//
// 	// So let SegmentAvailabilityEndTime equal the current time and compute
// 	// SegmentStart, which yields the start time that a segment would need to
// 	// have to have an availability end time equal to the current time.
// 	//
// 	// TODO: Use availabilityTimeOffset
// 	earliestAvailableTimestamp := currentPresentationTime - (2 * scaledSegmentDuration) - mpd.TimeShiftBufferDepth
// 	if earliestAvailableTimestamp < 0 {
// 		earliestAvailableTimestamp = 0
// 	}

// 	// Now round up to the nearest segment boundary, since the segment
// 	// corresponding to |earliestAvailableTimestamp| is not available.
// 	earliestAvailableSegmentStartTime := Math.ceil(earliestAvailableTimestamp/scaledSegmentDuration) * scaledSegmentDuration

// 	// Compute the segment start time of the latest available segment, i.e., the
// 	// segment that starts closest to the present but is available.
// 	//
// 	// Using the above formulas, let SegmentAvailabilityStartTime equal the
// 	// current time and compute SegmentStart, which yields the start time that
// 	// a segment would need to have to have an availability start time
// 	// equal to the current time.
// 	latestAvailableTimestamp := currentPresentationTime - scaledSegmentDuration
// 	if latestAvailableTimestamp < 0 {
// 		fmt.Printf("The first segment is not available yet!%s\r\n", period)
// 		return nil
// 	}

// 	// Now round down to the nearest segment boundary, since the segment
// 	// corresponding to |latestAvailableTimestamp| may not yet be available.
// 	latestAvailableSegmentStartTime := Math.floor(latestAvailableTimestamp/scaledSegmentDuration) * scaledSegmentDuration

// 	// Now compute the start time of the "best" available segment, by offsetting
// 	// by @suggestedPresentationDelay and @minBufferTime. Note that we subtract
// 	// by @minBufferTime to ensure that after playback begins we can buffer at
// 	// least @minBufferTime seconds worth of media content.
// 	bestAvailableTimestamp := latestAvailableSegmentStartTime - suggestedPresentationDelay - minBufferTime
// 	if bestAvailableTimestamp < 0 {
// 		fmt.Println("The first segment may not be available yet.")
// 		bestAvailableTimestamp = 0
// 		// Don't return; taking into account @suggestedPresentationDelay is only a
// 		// reccomendation. The first segment /might/ be available.
// 	}

// 	bestAvailableSegmentStartTime := Math.floor(bestAvailableTimestamp/scaledSegmentDuration) * scaledSegmentDuration

// 	// Now take the larger of |bestAvailableSegmentStartTime| and
// 	// |earliestAvailableSegmentStartTime|.
// 	var currentSegmentStartTime int

// 	if bestAvailableSegmentStartTime >= earliestAvailableSegmentStartTime {
// 		currentSegmentStartTime = bestAvailableSegmentStartTime
// 		fmt.Println("The best available segment is still available!")
// 	} else {
// 		// NOTE: @suggestedPresentationDelay + @minBufferTime is large compared to
// 		// @timeShiftBufferDepth, so we can't start as far back, for buffering, as
// 		// we'd like.
// 		currentSegmentStartTime = earliestAvailableSegmentStartTime
// 		fmt.Println("The best available segment is no longer available.")
// 	}

// 	earliestSegmentNumber := (earliestAvailableSegmentStartTime / scaledSegmentDuration) + 1
// 	assert(earliestSegmentNumber == Math.round(earliestSegmentNumber), "earliestSegmentNumber should be an integer.")

// 	currentSegmentNumber := (currentSegmentStartTime / scaledSegmentDuration) + 1
// 	assert(currentSegmentNumber == Math.round(currentSegmentNumber), "currentSegmentNumber should be an integer.")

// 	fmt.Printf("earliestSegmentNumber%s\r\n", earliestSegmentNumber)
// 	fmt.Printf("currentSegmentNumber%s\r\n", currentSegmentNumber)

// 	return NewPair(earliestSegmentNumber, currentSegmentNumber)
// }

/**
 * Generates an Initialization from a SegmentTemplate.
 *
 * @param {Representation} representation
 * @return {Initialization} An Initialization on success, null
 *     if no initialization URL template exists or an error occurred.
 */
func (mpdProcessor *MpdProcessor) generateInitialization(representation Representation) (Initialization, error) {
	assert(representation.SegmentTemplate != nil)

	initialization := Initialization{}

	segmentTemplate := representation.SegmentTemplate
	assert(segmentTemplate.InitializationUrlTemplate != "")
	if segmentTemplate.InitializationUrlTemplate == "" {
		return initialization, errors.New("segment template initialization url template is missing")
	}

	// $Number$ and $Time$ cannot be present in an initialization URL template.
	filledUrlTemplate := mpdProcessor.fillUrlTemplate(segmentTemplate.InitializationUrlTemplate, representation.Id, 0, representation.Bandwidth, 0)

	if filledUrlTemplate == "" {
		// An error has already been logged.
		return initialization, errors.New("could not fill initialization url template")
	}

	if representation.BaseUrl != nil && filledUrlTemplate != "" {
		initialization.Url = representation.BaseUrl.Url + filledUrlTemplate
	} else {
		initialization.Url = filledUrlTemplate
	}

	return initialization, nil
}

/**
 * Fills a SegmentTemplate URL template.
 *
 * @see ISO/IEC 23009-1:2014 section 5.3.9.4.4
 *
 * @param {string} urlTemplate
 * @param {?string} representationId
 * @param {?number} number
 * @param {?number} bandwidth
 * @param {?number} time
 * @return {string} A URL on success; null if the resulting URL contains
 *     illegal characters.
 */
func (mpdProcessor *MpdProcessor) fillUrlTemplate(urlTemplate string, representationId string, number int, bandwidth int, time int) string {
	/** @type {!Object.<string, ?number|?string>} */
	// fmt.Printf("In FillUrlTemplate, urlTemplate: %s, representationId: %s, number: %d, bandwidth: %d, time: %d\r\n", urlTemplate, representationId, number, bandwidth, time)
	valueTable := make(map[string]string)
	valueTable["$RepresentationID$"] = representationId
	valueTable["$Number$"] = strconv.Itoa(number)
	valueTable["Bandwidth"] = strconv.Itoa(bandwidth)
	valueTable["Time"] = strconv.Itoa(time)

	re := regexp.MustCompile("\\$(RepresentationID|Number|Bandwidth|Time)?(?:%0([0-9]+)d)?\\$")
	// var re = /\$(RepresentationID|Number|Bandwidth|Time)?(?:%0([0-9]+)d)?\$/g;
	url := re.ReplaceAllStringFunc(urlTemplate, func(match string) string {
		if match == "$$" {
			return "$"
		}

		value, ok := valueTable[match]
		assert(ok == true)

		// Note that |value| may be 0 or ''.
		if ok == false {
			fmt.Printf("URL template does not have an available substitution for identifier %s\r\n", match)
			return match
		}

		// if (match == "RepresentationID" && widthString) {
		//   fmt.Printf("URL template should not contain a width specifier for identifier %s\r\n", RepresentationID)
		//   widthString = undefined
		// }

		// Create padding string.
		// width := value
		// paddingSize := Math.max(0, width - valueString.length);
		// padding := (new Array(paddingSize + 1)).join('0');
		return value
		// return padding + valueString
	})

	// The URL might contain illegal characters (e.g., '%').
	return url
}

/**
 * Creates a SegmentMetadataInfo from either a RepresentationIndex or an
 * Initialization.
 *
 * @param {RepresentationIndex| Initialization} urlTypeObject
 * @return {SegmentMetadataInfo}
 */
func (mpdProcessor *MpdProcessor) createSegmentMetadataInfo(urlTypeObject Node) (SegmentMetadataInfo, error) {
	segmentMetadataInfo := NewSegmentMetadataInfo()

	if urlTypeObject == nil {
		return segmentMetadataInfo, errors.New("missing url type object")
	}

	var url string
	var r *Range

	switch obj := urlTypeObject.(type) {
	case *RepresentationIndex:
	case *Initialization:
		url = obj.Url
		r = obj.Range
	}

	segmentMetadataInfo.Url = url

	if r != nil {
		segmentMetadataInfo.StartByte = r.Begin
		segmentMetadataInfo.EndByte = r.End
	}

	return segmentMetadataInfo, nil
}
