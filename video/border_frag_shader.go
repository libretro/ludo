package video

var borderFragmentShader = `
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

uniform float border_width;
uniform vec4 color;
uniform vec2 size;

COMPAT_VARYING vec2 fragTexCoord;

void main() {
	float ratio = size.x / size.y;
	float maxX = 1.0-border_width/ratio;
	float minX = border_width/ratio;
	float maxY = 1.0-border_width;
	float minY = border_width;

	if (fragTexCoord.x < maxX && fragTexCoord.x > minX &&
			fragTexCoord.y < maxY && fragTexCoord.y > minY) {
		COMPAT_FRAGCOLOR = vec4(0,0,0,0);
	} else {
		COMPAT_FRAGCOLOR = color;
	}
}
` + "\x00"
