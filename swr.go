/*
2015  Sleepy Programmer <hunan@emsym.com>
*/
package gmf

/*

#cgo pkg-config: libswresample

#include "libswresample/swresample.h"
#include <libavcodec/avcodec.h>
#include <libavutil/frame.h>

int gmf_sw_resample(SwrContext* ctx, AVFrame*dstFrame, AVFrame*srcFrame){
	return swr_convert(ctx, dstFrame->data, dstFrame->nb_samples,
		(const uint8_t **)srcFrame->data, srcFrame->nb_samples);
}

*/
import "C"

type SwrCtx struct {
	swrCtx *C.struct_SwrContext
	cc     *CodecCtx
	inputSampleRate int
	CgoMemoryManage
}

func NewSwrCtx(options []*Option, cc *CodecCtx) *SwrCtx {
	this := &SwrCtx{swrCtx: C.swr_alloc(), cc: cc}

	for _, option := range options {
		if option.Key == "in_sample_rate" {
			this.inputSampleRate = option.Val.(int)
		}
		option.Set(this.swrCtx)
	}

	if int(C.swr_init(this.swrCtx)) < 0 {
		return nil
	}

	return this
}

func (this *SwrCtx) Free() {
	C.swr_free(&this.swrCtx)
}
func (this *SwrCtx) delay() int64 {

	return int64(C.swr_get_delay(this.swrCtx, C.int64_t(this.cc.avCodecCtx.sample_rate)))
}
func (this *SwrCtx) InputSampleRate() int {
	return this.inputSampleRate
}
func (this *SwrCtx) Convert(input *Frame) *Frame {
	if this.cc == nil {
		return nil
	}
	delaySamples := this.delay()
	srcRate := this.inputSampleRate
	dstRate := this.cc.SampleRate()
	dstSamples := int(RescaleRound(delaySamples+int64(input.NbSamples()),
		int64(dstRate), int64(srcRate), AV_ROUND_UP))
	channels := this.cc.Channels()
	format := this.cc.SampleFmt()
	dstFrame, _ := NewAudioFrame(format, channels, dstSamples)

	C.gmf_sw_resample(this.swrCtx, dstFrame.avFrame, input.avFrame)
	// frame := NewFrame()

	// dstNbSamples := C.av_rescale_rnd(C.swr_get_delay(this.swrCtx, this.cc.avCodecCtx.sample_rate)+input.avFrame.nb_samples, C.int64_t(this.cc.SampleRate()), this.cc.SampleRate(), C.AV_ROUND_UP)

	// if dstNbSamples > input.NbSamples() {
	// }

	return dstFrame
}
