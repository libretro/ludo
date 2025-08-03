package video

var uiFragmentShader = `
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

COMPAT_VARYING vec2 fragTexCoord;

uniform vec4 rect;
uniform float radius;
uniform vec4 color;
uniform vec4 borderColor;
uniform float borderWidth;
uniform float shadowSize;

float roundedRectSDF(vec2 p, vec2 size, float radius) {
    vec2 d = abs(p) - (size - vec2(radius));
    return length(max(d, 0.0)) + min(max(d.x, d.y), 0.0) - radius;
}

void main() {
    vec2 p = fragTexCoord - (rect.xy + rect.zw * 0.5);
    float dist = roundedRectSDF(p, rect.zw * 0.5, radius);
    float alpha = smoothstep(-shadowSize, 2.0, -dist);
    float borderDist = abs(dist) - borderWidth;
    vec4 colorFinal = mix(borderColor, color, smoothstep(0.0, 2.0, borderDist));
    COMPAT_FRAGCOLOR = vec4(colorFinal.rgb, alpha * colorFinal.a);
}
` + "\x00"
