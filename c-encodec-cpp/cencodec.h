#ifndef H_C_ENCODEC_CPP
#define H_C_ENCODEC_CPP

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct cencodec_compressed {
  int32_t* codes;
  int codes_len;
} cencodec_compressed;

typedef struct cencodec_decompressed {
  float* audio;
  int audio_len;
} cencodec_decompressed;

void* cencodec_load_model(const char* model_path, int n_gpu_layers);

void cencodec_set_target_bandwidth(void* ectx, int bandwidth);

int cencodec_reconstruct_audio(void* ectx, float* raw_audio, int raw_audio_len, int n_threads);

int cencodec_compress_audio(void* ectx, float* raw_audio, int raw_audio_len, int n_threads);

cencodec_compressed* cencodec_get_compress_codes(void* ectx);

void cencodec_compressed_free(cencodec_compressed* data);

int cencodec_decompress_audio(void* ectx, int32_t* codes, int codes_len, int n_threads);

cencodec_decompressed* cencodec_get_decompress_audio(void* ectx);

void cencodec_decompressed_free(cencodec_decompressed* data);

void cencodec_free(void* ectx);

#ifdef __cplusplus
};
#endif /* __cplusplus */

#endif // H_C_ENCODEC_CPP
