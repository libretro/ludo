package video

var roundedFragmentShader = `
#version 330

uniform vec4 color;
uniform float radius;
uniform vec2 size;

in vec2 fragTexCoord;
out vec4 outputColor;

float udRoundBox(vec2 p, vec2 b, float r) {
  return length(max(abs(p)-b+r,0.0))-r;
}

void main() {
	float ratio = size.x / size.y;
	vec2 halfRes = vec2(0.5*ratio, 0.5);
	float b = udRoundBox(fragTexCoord*vec2(ratio,1.0) - halfRes, halfRes, min(halfRes.x,halfRes.y)*radius);
	outputColor = vec4(color.r, color.g, color.b, min(color.a, 1-smoothstep(0.00001,0.001,b)));
}
` + "\x00"
