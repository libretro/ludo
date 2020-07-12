package video

var circleFragmentShader = `
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

COMPAT_VARYING vec2 fragTexCoord;

float circle(vec2 _st, float _radius) {
  vec2 dist = _st - vec2(0.5);
  return 1.-smoothstep(_radius-(_radius*0.05), _radius+(_radius*0.05), dot(dist,dist)*4.0);
}

void main() {
	COMPAT_FRAGCOLOR = vec4(color.rgb, circle(fragTexCoord.xy, 0.125)*color.a);
}
` + "\x00"
