package video

// source of the shader to draw circles
var circleFragmentShader = `
#version 330

uniform sampler2D tex;
uniform vec4 texColor;

in vec2 fragTexCoord;
out vec4 outputColor;

float circle(in vec2 _st, in float _radius) {
  vec2 dist = _st - vec2(0.5);
  return 1.-smoothstep(_radius-(_radius*0.01), _radius+(_radius*0.01), dot(dist,dist)*4.0);
}

void main() {
  vec4 color = texture(tex, fragTexCoord);
	outputColor = texColor * color;
	outputColor = vec4(outputColor.rgb, circle(fragTexCoord.xy, 1.0));
}
` + "\x00"
