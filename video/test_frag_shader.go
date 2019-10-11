package video

var testFragmentShader = `
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

uniform vec4 box;
uniform vec4 color;
uniform float radius;
uniform vec2 resolution;

COMPAT_VARYING vec2 fragTexCoord;

float rec(vec2 uv, vec2 tl, vec2 br)
{
	vec2 d = max(tl-uv, uv-br);
	return length(max(vec2(0.0), d)) + min(0.0, max(d.x, d.y));
}

void main() {
	vec2 st = gl_FragCoord.xy/resolution;
	st -= 0.5;
	st *= vec2(resolution.x/resolution.y, -1.0);
	st += vec2(0.5*resolution.x/resolution.y, 0.5);

	float d = rec(st, box.xy/resolution.y, box.zw/resolution.y);
	float m = 1.0 - d/radius;
	COMPAT_FRAGCOLOR = vec4(color.rgb, m*color.a);
}
` + "\x00"
