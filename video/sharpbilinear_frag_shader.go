package video

var sharpBilinearFragmentShader = `
#if __VERSION__ >= 130
#define COMPAT_VARYING in
#define COMPAT_ATTRIBUTE in
#define COMPAT_TEXTURE texture
#define COMPAT_FRAGCOLOR FragColor
out vec4 COMPAT_FRAGCOLOR;
#else
#define COMPAT_VARYING varying
#define COMPAT_ATTRIBUTE attribute
#define COMPAT_TEXTURE texture2D
#define COMPAT_FRAGCOLOR gl_FragColor
#endif

uniform vec2 OutputSize;
uniform vec2 TextureSize;
uniform vec2 InputSize;
uniform sampler2D Texture;
COMPAT_VARYING vec2 fragTexCoord;

// fragment compatibility #defines
#define SourceSize vec4(TextureSize, 1.0 / TextureSize) //either TextureSize or InputSize
#define outsize vec4(OutputSize, 1.0 / OutputSize)

void main() {
  vec2 texel = fragTexCoord * SourceSize.xy;
  vec2 scale = max(floor(outsize.xy / InputSize.xy), vec2(1.0, 1.0));

  vec2 texel_floored = floor(texel);
  vec2 s = fract(texel);
  vec2 region_range = 0.5 - 0.5 / scale;

  vec2 center_dist = s - 0.5;
  vec2 f = (center_dist - clamp(center_dist, -region_range, region_range)) * scale + 0.5;

  vec2 mod_texel = texel_floored + f;

  COMPAT_FRAGCOLOR = vec4(COMPAT_TEXTURE(Texture, mod_texel / SourceSize.xy).rgb, 1.0);
}
` + "\x00"
