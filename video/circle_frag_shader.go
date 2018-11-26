package video

// source of the shader to draw circles
var circleFragmentShader = `
#version 120

uniform sampler2D tex;
uniform vec4 color;

varying vec2 fragTexCoord;

float circle(vec2 _st, float _radius) {
  vec2 dist = _st - vec2(0.5);
  return 1.-smoothstep(_radius-(_radius*0.05), _radius+(_radius*0.05), dot(dist,dist)*4.0);
}

void main() {
	gl_FragColor = vec4(color.rgb, circle(fragTexCoord.xy, 0.125));
}
` + "\x00"
