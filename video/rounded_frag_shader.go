package video

var roundedFragmentShader = `
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

uniform vec4 color;
uniform float radius;
uniform vec2 size;

COMPAT_VARYING vec2 fragTexCoord;

float udRoundBox(vec2 p, vec2 b, float r) {
  return length(max(abs(p)-b+r,0.0))-r;
}

void main() {
	float ratio = size.x / size.y;
	vec2 halfRes = vec2(0.5*ratio, 0.5);
	float b = udRoundBox(fragTexCoord*vec2(ratio,1.0) - halfRes, halfRes, min(halfRes.x,halfRes.y)*radius);
	COMPAT_FRAGCOLOR = vec4(color.r, color.g, color.b, min(color.a, 1.0-smoothstep(0.00001,0.001,b)));
}
` + "\x00"
