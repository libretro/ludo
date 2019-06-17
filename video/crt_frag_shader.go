package video

var zfastCRTFragmentShader = `
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

#define BLURSCALEX 0.45
#define LOWLUMSCAN 5.0
#define HILUMSCAN 10.0
#define BRIGHTBOOST 1.25
#define MASK_DARK 0.25
#define MASK_FADE 0.8

#define Source Texture
#define vTexCoord fragTexCoord.xy

void main() {
  float maskFade = 0.3333*MASK_FADE;
  vec2 invDims = 1.0/TextureSize.xy;

  vec2 p = vTexCoord * TextureSize;
  vec2 i = floor(p) + 0.50;
  vec2 f = p - i;

  p = (i + 4.0*f*f*f)*invDims;
  p.x = mix( p.x , vTexCoord.x, BLURSCALEX);
  float Y = f.y*f.y;
  float YY = Y*Y;

  float whichmask = fract( gl_FragCoord.x*-0.4999);
  float mask = 1.0 + float(whichmask < 0.5) * -MASK_DARK;

  vec3 colour = COMPAT_TEXTURE(Source, p).rgb;

  float scanLineWeight = (BRIGHTBOOST - LOWLUMSCAN*(Y - 2.05*YY));
  float scanLineWeightB = 1.0 - HILUMSCAN*(YY-2.8*YY*Y);

  COMPAT_FRAGCOLOR.rgb = colour.rgb*mix(scanLineWeight*mask, scanLineWeightB, dot(colour.rgb,vec3(maskFade)));
}
` + "\x00"
