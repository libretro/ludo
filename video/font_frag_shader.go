package video

var fontFragmentShader = `
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

COMPAT_VARYING vec2 fragTexCoord;

void main() {
	vec4 sampled = vec4(1.0, 1.0, 1.0, COMPAT_TEXTURE(Texture, fragTexCoord).r);
	COMPAT_FRAGCOLOR = min(color, vec4(1.0, 1.0, 1.0, 1.0)) * sampled;
}
` + "\x00"
