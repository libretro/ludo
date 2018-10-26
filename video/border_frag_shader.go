package video

// source of the shader to draw circles
var borderFragmentShader = `
#version 330

uniform float border_width;
uniform vec4 color;
uniform vec2 size;

in vec2 fragTexCoord;
out vec4 outputColor;

void main() {
	float ratio = size.x / size.y;
	float maxX = 1.0-border_width/ratio;
	float minX = border_width/ratio;
	float maxY = 1.0-border_width;
	float minY = border_width;

	if (fragTexCoord.x < maxX && fragTexCoord.x > minX &&
			fragTexCoord.y < maxY && fragTexCoord.y > minY) {
		outputColor = vec4(0,0,0,0);
	} else {
		outputColor = color;
	}
}
` + "\x00"
