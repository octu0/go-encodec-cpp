#include <cstring>
#include <cstdlib>
#include <string>
#include "encodec.h"
#include "cencodec.h"

extern "C" {

void* cencodec_load_model(const char* model_path, int n_gpu_layers){
  struct encodec_context* ctx = encodec_load_model(std::string(model_path), n_gpu_layers);
  return ctx;
}

void cencodec_set_target_bandwidth(void* ectx, int bandwidth) {
  encodec_set_target_bandwidth((struct encodec_context*) ectx, bandwidth);
}

int cencodec_reconstruct_audio(void* ectx, float* raw_audio, int raw_audio_len, int n_threads) {
  struct encodec_context* ctx = (struct encodec_context*) ectx;
  std::vector<float> audio_arr;
  audio_arr.resize(raw_audio_len);
  memcpy(audio_arr.data(), raw_audio, raw_audio_len * sizeof(float));
  bool ok = encodec_reconstruct_audio(ctx, audio_arr, n_threads);
  return ok ? 0 : 1;
}

int cencodec_compress_audio(void* ectx, float* raw_audio, int raw_audio_len, int n_threads) {
  struct encodec_context* ctx = (struct encodec_context*) ectx;
  std::vector<float> audio_arr;
  audio_arr.resize(raw_audio_len);
  memcpy(audio_arr.data(), raw_audio, raw_audio_len * sizeof(float));
  bool ok = encodec_compress_audio(ctx, audio_arr, n_threads);
  return ok ? 0 : 1;
}

cencodec_compressed* cencodec_get_compress_codes(void* ectx) {
  struct encodec_context* ctx = (struct encodec_context*) ectx;
  cencodec_compressed* data = (cencodec_compressed*) malloc(sizeof(cencodec_compressed));
  if(nullptr == data) {
    return nullptr;
  }
  memset(data, 0, sizeof(cencodec_compressed));

  data->codes_len = ctx->out_codes.size();
  data->codes = (int32_t*) malloc(data->codes_len * sizeof(int32_t));
  if(nullptr == data->codes) {
    free(data);
    return nullptr;
  }
  memcpy(data->codes, ctx->out_codes.data(), data->codes_len * sizeof(int32_t));
  return data;
}

void cencodec_compressed_free(cencodec_compressed* data) {
  if(nullptr != data) {
    free(data->codes);
  }
  free(data);
}

int cencodec_decompress_audio(void* ectx, int32_t* codes, int codes_len, int n_threads) {
  struct encodec_context* ctx = (struct encodec_context*) ectx;
  std::vector<int32_t> codes_arr;
  codes_arr.resize(codes_len);
  memcpy(codes_arr.data(), codes, codes_len * sizeof(int32_t));
  bool ok = encodec_decompress_audio(ctx, codes_arr, n_threads);
  return ok ? 0 : 1;
}

cencodec_decompressed* cencodec_get_decompress_audio(void* ectx) {
  struct encodec_context* ctx = (struct encodec_context*) ectx;
  cencodec_decompressed* data = (cencodec_decompressed*) malloc(sizeof(cencodec_decompressed));
  if(nullptr == data){
    return nullptr;
  }
  memset(data, 0, sizeof(cencodec_decompressed));

  data->audio_len = ctx->out_audio.size();
  data->audio = (float*) malloc(data->audio_len * sizeof(float));
  if(nullptr == data->audio){
    free(data);
    return nullptr;
  }
  memcpy(data->audio, ctx->out_audio.data(), data->audio_len * sizeof(float));
  return data;
}

void cencodec_decompressed_free(cencodec_decompressed* data) {
  if(nullptr != data) {
    free(data->audio);
  }
  free(data);
}

void cencodec_free(void* ectx) {
  encodec_free((struct encodec_context*) ectx);
}

}
