package video

var zfastLCDFragmentShader = `
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

#define BORDERMULT 14.0
#define GBAGAMMA 1.0

void main() {
  vec2 texcoordInPixels = fragTexCoord.xy * TextureSize.xy;
  vec2 centerCoord = floor(texcoordInPixels.xy)+vec2(0.5,0.5);
  vec2 distFromCenter = abs(centerCoord - texcoordInPixels);
  vec2 invSize = 1.0/TextureSize.xy;

  float Y = max(distFromCenter.x,(distFromCenter.y));

  Y=Y*Y;
  float YY = Y*Y;
  float YYY = YY*Y;

  float LineWeight = YY - 2.7*YYY;
  LineWeight = 1.0 - BORDERMULT*LineWeight;

  vec3 colour = COMPAT_TEXTURE(Texture, invSize*centerCoord).rgb*LineWeight;

  if (GBAGAMMA > 0.5)
    colour.rgb*=0.6+0.4*(colour.rgb); //fake gamma because the pi is too slow!
    
  COMPAT_FRAGCOLOR = vec4(colour.rgb , 1.0);
}
` + "\x00"
