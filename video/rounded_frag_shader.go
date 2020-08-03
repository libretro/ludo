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

uniform sampler2D Texture;
uniform vec4 color;
uniform float radius;
uniform vec2 size;

COMPAT_VARYING vec2 fragTexCoord;

float udRoundBox(vec2 p, vec2 b, float r) {
  return length(max(abs(p)-b+r,0.0))-r;
}

vec4 demultiply(vec4 c) {
  return vec4(c.rgb/c.a, c.a);
}

void main() {
	float ratio = size.x / size.y;
	vec2 halfRes = vec2(0.5*ratio, 0.5);
	float b = udRoundBox(fragTexCoord*vec2(ratio,1.0) - halfRes, halfRes, min(halfRes.x,halfRes.y)*radius);
	vec4 c = min(color, vec4(1.0, 1.0, 1.0, 1.0));
	COMPAT_FRAGCOLOR = demultiply(COMPAT_TEXTURE(Texture, fragTexCoord)) * vec4(c.r, c.g, c.b, c.a * (1.0-smoothstep(0.0001,0.001,b)));
}
` + "\x00"
