package video

var uiVertexShader = `
#if __VERSION__ >= 130
#define COMPAT_VARYING out
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

COMPAT_ATTRIBUTE vec2 vert;
COMPAT_ATTRIBUTE vec2 vertTexCoord;
COMPAT_VARYING vec2 fragTexCoord;

uniform vec4 rect;
uniform vec2 resolution;
uniform float shadowSize;

void main() {
    vec2 paddedSize = rect.zw*20.0 + vec2(shadowSize * 2.0);
    vec2 pos = rect.xy + vert * paddedSize - vec2(shadowSize);
    fragTexCoord = pos;
    vec2 ndc = (pos / resolution) * 2.0 - 1.0;
    gl_Position = vec4(ndc.x, -ndc.y, 0.0, 1.0);
}
` + "\x00"
