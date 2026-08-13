package fairy

import (
	"encoding/base64"

	"github.com/kirinyoku/fairy/store"
)

// Attribute represents the elemental attribute of an agent.
type Attribute string

const (
	// AttributePhysical represents the agent's Physical attribute.
	AttributePhysical Attribute = "Physical"
	// AttributeHonedEdge represents the agent's Honed Edge attribute.
	AttributeHonedEdge Attribute = "HonedEdge"
	// AttributeFire represents the agent's Fire attribute.
	AttributeFire Attribute = "Fire"
	// AttributeIce represents the agent's Ice attribute.
	AttributeIce Attribute = "Ice"
	// AttributeFrost represents the agent's Frost attribute.
	AttributeFrost Attribute = "Frost"
	// AttributeElectric represents the agent's Electric attribute.
	AttributeElectric Attribute = "Electric"
	// AttributeEther represents the agent's Ether attribute.
	AttributeEther Attribute = "Ether"
	// AttributeAuricInk represents the agent's Auric Ink attribute.
	AttributeAuricInk Attribute = "AuricInk"
	// AttributeWind represents the agent's Wind attribute.
	AttributeWind Attribute = "Wind"
	// AttributeLumiflux represents the agent's Lumiflux attribute.
	AttributeLumiflux Attribute = "Lumiflux"
)

// Specialty represents the combat role or class of an agent.
type Specialty string

const (
	// SpecialtyAttack represents the Attack combat role.
	SpecialtyAttack Specialty = "Attack"
	// SpecialtyStun represents the Stun combat role.
	SpecialtyStun Specialty = "Stun"
	// SpecialtyAnomaly represents the Anomaly combat role.
	SpecialtyAnomaly Specialty = "Anomaly"
	// SpecialtySupport represents the Support combat role.
	SpecialtySupport Specialty = "Support"
	// SpecialtyDefense represents the Defense combat role.
	SpecialtyDefense Specialty = "Defense"
	// SpecialtyRupture represents the Rupture combat role.
	SpecialtyRupture Specialty = "Rupture"
)

// Rarity represents the rarity tier of agents and equipment.
type Rarity string

const (
	// RarityS represents the S-rank tier.
	RarityS Rarity = "S"
	// RarityA represents the A-rank tier.
	RarityA Rarity = "A"
	// RarityB represents the B-rank tier.
	RarityB Rarity = "B"
)

var attributeSVGMap = map[Attribute]string{
	AttributeAuricInk:  `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_AuricEther" viewBox="0 0 14 14"><defs><linearGradient id="zzz-AddedDamageRatio_AuricEther__a"><stop offset="0" style="stop-color:#f3d299;stop-opacity:1"/><stop offset=".5" style="stop-color:#b7822c;stop-opacity:1"/><stop offset="1" style="stop-color:#ffcd77;stop-opacity:1"/></linearGradient><linearGradient href="#zzz-AddedDamageRatio_AuricEther__a" id="zzz-AddedDamageRatio_AuricEther__b" x1="15.915" x2="22.39" y1="2.392" y2="8.867" gradientTransform="translate(-12.386 1.207)" gradientUnits="userSpaceOnUse"/></defs><path d="M7.161 0S5.305 2.79 4.128 3.967c-.254.253-.59.54-.948.831l.076.063c2.105 1.177 3.63 2.942 5.05 4.666l-1.253 1.621L5.95 9.585l.784-.833c-.02-.198-.513-.735-1.148-1.326-.508.47-1.474 1.312-2.384 1.795.35.284.677.564.925.812C5.305 11.21 7.033 14 7.033 14c1.162-1.719 2.064-3.582 3.957-4.78-.533-.283-1.803-1.212-2.619-2.22l-1.21-1.322-1.21-1.287 1.268-1.588L8.294 4.32l-.91.955 1.205 1.308c.536-.69 1.92-1.806 2.41-1.904C9.69 3.6 8.389 2.143 7.16 0m-4.1 5.344S1.154 7 .652 7c.5 0 2.41 1.656 2.41 1.656l1.369-.898-1.04-.595v-.327l1.04-.594Zm7.878 0-1.37.898 1.04.594v.327l-1.04.595 1.37.898S12.904 7 13.348 7c-.444 0-2.41-1.656-2.41-1.656" style="opacity:1;fill:url(#zzz-AddedDamageRatio_AuricEther__b);stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`,
	AttributeElectric:  `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Elec" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Elec__a" id="zzz-AddedDamageRatio_Elec__b" x1="12.046" x2="12.046" y1="334.813" y2="349.957" gradientTransform="translate(-4.136 -309.526)scale(.92447)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Elec__a"><stop offset="0" style="stop-color:#0075ff;stop-opacity:1"/><stop offset="1" style="stop-color:#3decff;stop-opacity:1"/></linearGradient></defs><path d="m9.822.624-5.237.624-1.573 5.236 2.143-.054-1.736 6.946 6.783-8.628-2.768.217ZM1.6 5.779 0 7.705l.678 2.578.652-.57.949 2.496.597-4.531-1.167.814Zm9.17 1.004L8.14 9.225l1.166.298-2.7 3.392 6.065-2.768-1.764-.786L14 7Z" style="fill:url(#zzz-AddedDamageRatio_Elec__b);stroke-width:0.924468;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`,
	AttributeEther:     `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Ether" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Ether__a" id="zzz-AddedDamageRatio_Ether__b" x1="16.461" x2="7.36" y1="354.589" y2="367.587" gradientTransform="translate(-2.93 -294.054)scale(.83374)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Ether__a"><stop offset="0" style="stop-color:#ff0a1a;stop-opacity:1"/><stop offset=".171" style="stop-color:#ff0626;stop-opacity:1"/><stop offset=".5" style="stop-color:#b338dd;stop-opacity:1"/><stop offset=".85" style="stop-color:#2a6bea;stop-opacity:1"/><stop offset="1" style="stop-color:#2a6bea;stop-opacity:1"/></linearGradient></defs><path d="M6.52 0 5.036 3.715a1.68 1.68 0 0 1-.935.936L.385 6.135l3.716 1.483a1.47 1.47 0 0 1 .866.96L6.52 14l1.553-5.422a1.47 1.47 0 0 1 .866-.96l3.716-1.483L8.94 4.65a1.68 1.68 0 0 1-.935-.936Zm4.568 7.83-.57 1.403a1 1 0 0 1-.554.553l-1.403.57 1.403.57a1 1 0 0 1 .554.554l.57 1.404.57-1.404a1 1 0 0 1 .554-.553l1.403-.57-1.403-.57a1 1 0 0 1-.554-.554Z" style="fill:url(#zzz-AddedDamageRatio_Ether__b);stroke-width:0.833738;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`,
	AttributeFire:      `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Fire" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Fire__a" id="zzz-AddedDamageRatio_Fire__b" x1="12.182" x2="12.182" y1="302.124" y2="315.03" gradientTransform="translate(-4.909 -294.668)scale(.9776)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Fire__a"><stop offset="0" style="stop-color:#ea1503;stop-opacity:1"/><stop offset="1" style="stop-color:#f3741a;stop-opacity:1"/></linearGradient></defs><path d="M1.972 11.944c1.112 1.072 2.998 1.314 4.368 2.027.007-1.742-1.044-3.54-3.158-3.855 2.096-.47 3.1-1.295 3.665-3.693.572 2.723 1.85 3.219 3.91 3.678-2.092.376-3.42 1.696-3.159 3.899 1.55-.806 3.624-1.29 4.76-2.677 2.117-2.58.435-6.764-2.47-7.9.949 2.315-1.854 2.066-2.105-.375C7.666 1.909 8.886.993 9.77.47 7.84-.68 4.008.404 4.242 2.366c.191 1.6 1.52 3.367-.155 3.766-1.715.409-1.2-2.104-1.2-2.104C.665 5.06-.254 9.8 1.973 11.945" style="fill:url(#zzz-AddedDamageRatio_Fire__b);fill-opacity:1;stroke:none;stroke-width:0.977604"/></svg>`,
	AttributeFrost:     `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_FireFrost" viewBox="0 0 14 14"><defs><linearGradient id="zzz-AddedDamageRatio_FireFrost__a"><stop offset="0" style="stop-color:#28fbff;stop-opacity:1"/><stop offset=".228" style="stop-color:#21f8f8;stop-opacity:1"/><stop offset=".5" style="stop-color:#3cd6f1;stop-opacity:1"/><stop offset=".792" style="stop-color:#73a2ff;stop-opacity:1"/><stop offset="1" style="stop-color:#75a3ff;stop-opacity:1"/></linearGradient><linearGradient href="#zzz-AddedDamageRatio_FireFrost__a" id="zzz-AddedDamageRatio_FireFrost__b" x1="10.74" x2="3.26" y1="1.242" y2="12.758" gradientUnits="userSpaceOnUse"/></defs><path d="M7.021 0c-.111 0-.5 1.435-1.35 2.226-.034.032-.003.122-.003.14.547.278 1.03.663 1.35 1.143.319-.575.773-.936 1.31-1.174 0 0 .037-.114 0-.146C7.519 1.482 7.133 0 7.022 0M.134 2.047c.109.24.306.355.481.547.209.227.39.518.537.79.382.708.54 1.53 1.01 2.187.836 1.163 2.65 1.723 3.987 1.094l-.42-.492-.613-.845c.616.25 1.324.927 1.878.921.634-.006 1.227-.68 1.829-.921a2.13 2.13 0 0 1-.608.898c-.105.09-.477.395-.477.395s.12.202.963.224c.779.045 1.805-.01 2.43-.55.87-.75 1.08-1.707 1.596-2.669.177-.33.365-.675.614-.958.172-.196.384-.342.525-.56-.302-.11-.722-.008-1.033.069-.777.192-1.392.431-2.066.876-.186.122-.687.46-.763.027-.052-.302.423-.585.52-.851-1.685.168-3.087.431-3.528 2.286-.218-1.452-2.124-2.354-3.58-2.286.155.245.858.66.516 1.013-.24.247-.735-.235-.942-.355-.881-.508-1.826-.836-2.856-.84m4.614 5.052c-.642.01-1.289.143-1.758.447-.338.218-.686.612-.912.941-.425.622-.57 1.405-.926 2.066a4 4 0 0 1-.598.85c-.163.17-.336.261-.42.487a5.53 5.53 0 0 0 2.917-.87c.183-.118.684-.53.872-.278.233.314-.31.67-.386.966 1.381-.07 3.041-.512 3.341-2.066.368.25.434.765.735 1.094.551.601 1.362.749 2.121.874.267.044.598.163.85.098-.136-.212-.74-.551-.618-.846.159-.386.594-.057.801.068a7.1 7.1 0 0 0 3.038 1.02c-.175-.478-.7-.874-.957-1.336-.426-.763-.596-1.67-1.134-2.37-.942-1.226-2.606-1.136-3.99-.99l.6.57.499.846-1.276-.743-.593-.223-.622.282-1.276.684.5-.728.593-.67a5 5 0 0 0-1.4-.173m2.204 3.392c-.319.575-.773.936-1.31 1.174 0 0-.038.114-.001.146C6.453 12.518 6.837 14 6.949 14c.11 0 .499-1.435 1.35-2.226.034-.032.003-.122.003-.14-.548-.278-1.03-.663-1.35-1.143" style="fill:url(#zzz-AddedDamageRatio_FireFrost__b);stroke:none;stroke-width:0.229848"/></svg>`,
	AttributeIce:       `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Ice" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Ice__a" id="zzz-AddedDamageRatio_Ice__b" x1="12.046" x2="12.046" y1="318.508" y2="331.879" gradientTransform="translate(-3.923 -287.884)scale(.9068)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Ice__a"><stop offset="0" style="stop-color:#04c2c8;stop-opacity:1"/><stop offset="1" style="stop-color:#83f4f0;stop-opacity:1"/></linearGradient></defs><path d="M7 0 5.166 3.824.938 3.5 3.332 7 .938 10.5l4.228-.323L7 14l1.834-3.823 4.228.323L10.668 7l2.394-3.5-4.228.324ZM5.18 5.06a.1.1 0 0 1 .06.015L7 5.99l1.76-.914a.123.123 0 0 1 .166.167l-.914 1.76.914 1.759a.123.123 0 0 1-.166.166L7 8.013l-1.76.914a.123.123 0 0 1-.166-.166L5.988 7l-.914-1.76a.123.123 0 0 1 .106-.18" style="fill:url(#zzz-AddedDamageRatio_Ice__b);stroke-width:0.906792;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`,
	AttributeLumiflux:  `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Lumen" viewBox="0 0 58 58"><path fill="#fe9ac9" d="M33 29c0 2.21-1.79 4-4 4s-4-1.79-4-4 1.79-4 4-4 4 1.79 4 4m25 0-15.24-5.64L33.47 33a6 6 0 0 1-4.48 2l16.26 10.26-4.54-9.87L57.99 29Zm-42.76 5.64L24.53 25c1.1-1.22 2.69-2 4.48-2L12.74 12.74l4.54 9.87L0 29zm17.75-10.12A5.98 5.98 0 0 1 35 29l10.26-16.26-9.87 4.54L29 0l-5.64 15.24L33 24.53h-.01Zm-7.98 8.96A5.98 5.98 0 0 1 23 29L12.74 45.26l9.87-4.54L29 58l5.64-15.24L25 33.47h.01Z" class="cls-1" data-name="Layer 15"/></svg>`,
	AttributePhysical:  `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Physics" viewBox="0 0 14 14"><defs><linearGradient href="#zzz-AddedDamageRatio_Physics__a" id="zzz-AddedDamageRatio_Physics__b" x1="12.046" x2="12.046" y1="278.603" y2="299.007" gradientTransform="translate(-1.265 -191.16)scale(.68614)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_Physics__a"><stop offset="0" style="stop-color:#e78801;stop-opacity:1"/><stop offset="1" style="stop-color:#efd400;stop-opacity:1"/></linearGradient></defs><path d="M9.013.217c-.87 1.096-2.116 2.47-3.595 2.36-.822-.06-2.224-.78-2.348-.663-.119.148.464 1.196.444 1.806C3.442 5.958 0 6.851 0 6.971c0 .127 1.243.482 1.805.865 1.862 1.266.612 2.894.182 4.554 1.503-.514 2.918-1.891 4.504-.761.8.57 1.173 1.273 1.817 2.186.755-2.82.996-3.654 2.71-4.18.897-.277 2.053-.178 2.982-.178-.45-.435-.988-1.001-1.404-1.583-1.44-2.013.221-3.08 1.311-4.696-1.88.059-3.876.534-4.562-1.806-.106-.36-.128-.776-.14-1.148 0 0-.058-.037-.092-.039-.034-.001-.1.034-.1.034zm-.986 3c.213 1.86 1.299 1.586 2.449 1.587 0 0-.813.657-.813 1.625s.813 1.626.813 1.626c-1.877-.12-2.618.594-2.752 2.273C6.78 8.834 6.3 8.476 4.515 9.319c.848-1.776.097-2.214-1.344-2.606 1.176-.488 2.515-1.09 1.886-2.632 1.565.642 2.184.209 2.97-.863" style="fill:url(#zzz-AddedDamageRatio_Physics__b);stroke:none;stroke-width:0.686135"/></svg>`,
	AttributeWind:      `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_Wind" viewBox="0 0 14 14"><defs><linearGradient id="zzz-AddedDamageRatio_Wind__a"><stop offset="0" style="stop-color:#61a3ff;stop-opacity:1"/><stop offset="1" style="stop-color:#97e3fa;stop-opacity:1"/></linearGradient><linearGradient href="#zzz-AddedDamageRatio_Wind__a" id="zzz-AddedDamageRatio_Wind__b" x1="1862.372" x2="1908.889" y1="1828.644" y2="1900.273" gradientUnits="userSpaceOnUse"/><linearGradient href="#zzz-AddedDamageRatio_Wind__a" id="zzz-AddedDamageRatio_Wind__c" x1="1862.372" x2="1908.889" y1="1828.644" y2="1900.273" gradientUnits="userSpaceOnUse"/></defs><g style="fill:url(#zzz-AddedDamageRatio_Wind__b)"><path fill="#FFF" d="M1870.45 1829q10.4-.95 17.85.6 6.65 1.4 10.4 4.6-12.85-5.45-27.25-1.45-15.65 4.35-24.5 18.15-5.25 8.15-3.65 18.25 2-7.9 6.4-13.65 4.9-6.45 12.35-9.65-23.2 18.7-11.3 39.95 7.95 14.1 25.85 16.75 16.85 2.45 29.45-6.4-13.55 1.95-22.2.05-10.3-2.25-16.55-10.35 9.2 7.1 22 6.35 11.7-.7 22.2-7.5 10.55-6.85 14.7-16.65 4.55-10.75-.6-21.1-.789 17.455-14.05 25 9.75-10.45 9.55-20.85-.2-9.35-8.05-16.3-7.6-6.65-19.1-8.4-12.1-1.9-23.5 2.6m26.35 43.15q-7.2 5.35-15.8 5.35-8.55 0-13.55-5.35-4.9-5.4-3.25-12.9 1.65-7.55 8.95-13 7.25-5.3 15.85-5.3 8.55 0 13.45 5.3 5 5.45 3.35 13-1.65 7.5-9 12.9" style="fill:url(#zzz-AddedDamageRatio_Wind__c)" transform="translate(-299.382 -295.987)scale(.16249)"/></g></svg>`,
	AttributeHonedEdge: `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="SvgIcon AddedDamageRatio_ZhenZhenAssault" viewBox="0 0 24 24"><defs><radialGradient href="#zzz-AddedDamageRatio_ZhenZhenAssault__a" id="zzz-AddedDamageRatio_ZhenZhenAssault__b" cx="10.919" cy="8.958" r="10.566" fx="10.919" fy="8.958" gradientTransform="matrix(.89188 -.48417 .5089 .93742 -3.489 6.33)" gradientUnits="userSpaceOnUse"/><linearGradient id="zzz-AddedDamageRatio_ZhenZhenAssault__a"><stop offset="0" style="stop-color:#dfddf0;stop-opacity:1"/><stop offset=".682" style="stop-color:#869eee;stop-opacity:1"/><stop offset="1" style="stop-color:#8177b9;stop-opacity:1"/></linearGradient></defs><path d="M11.394 2.224c-.81 1.422-1.723 2.978-1.98 4.516.59-.25 1.076-.414 1.473-.36-.004.743.1 1.306.414 1.698 0 0 .511-1.042 1.71-1.936 1.2-.893 2.31-1.267 2.31-1.267-3.453.111-3.927-2.65-3.927-2.65m8.199 2.584c-1.866.459-3.645.722-5.32 1.74.575.297.925.703 1.092 1.08-.744.593-1.164 1.006-1.38 1.585 0 0 1.082-.325 2.665-.189 1.582.137 2.977 1.067 2.977 1.067-2.227-3.634-.034-5.283-.034-5.283M8.87 5.19c-2.76 2.026-4.864.855-4.864.855.534 1.357 1.097 2.86 2.158 3.912.257-.54.502-.963.856-1.2.607.428 1.152.68 1.719.693 0 0-.429-1.048-.42-2.047S8.87 5.19 8.87 5.19m2.463 3.07-.551 1.92-1.89-.608 1.18 1.565-1.696 1.13 2 .02-.22 2.01 1.313-1.529 1.53 1.29-.476-1.843 2.104-.237-1.879-.933 1.095-1.714-1.87.808zM4.427 9.63c-.192 3.398-2.754 4.364-2.754 4.364 1.592.39 3.335.845 5.037.634-.295-.527-.494-.974-.448-1.397.819-.224 1.436-.505 1.856-.946 0 0-1.29-.257-2.193-.926-.903-.67-1.498-1.73-1.498-1.73m13.206.964a2 2 0 0 0-.146.01c.098.373.038.747-.04 1.087-1.57-.04-2.038.052-2.589.288 0 0 1.116.62 1.999 1.893.882 1.272.9 3.253.9 3.253 1.384-4.238 5.049-3.667 5.049-3.667-1.55-1.083-3.64-2.905-5.173-2.864m-4.575 3.643s.171 1.267-.318 2.736c-.49 1.47-2.066 2.67-2.066 2.67 4.223-1.433 6.175 1.996 6.175 1.996-.063-1.95-.038-5.806-.856-6.227-.515.459-.774.7-1.453.629-.243-.908-.963-1.504-1.482-1.804m-3.006.355s-1.083.9-2.347 1.175-2.64.13-2.64.13c2.671 2.462 1.633 5.446 1.633 5.446a13.7 13.7 0 0 0 4.26-4.074c-.64-.121-1.043-.074-1.36-.397.409-.702.524-1.644.454-2.28" style="fill:url(#zzz-AddedDamageRatio_ZhenZhenAssault__b);stroke-width:1.67595;stroke-linecap:square;stroke-dashoffset:3.77952"/></svg>`,
}

// BaseAttribute returns the core elemental attribute that this attribute deals damage as.
// For example, AuricInk deals Ether DMG, HonedEdge deals Physical DMG, and Frost deals Ice DMG.
func (a Attribute) BaseAttribute() Attribute {
	switch a {
	case AttributeAuricInk:
		return AttributeEther
	case AttributeHonedEdge:
		return AttributePhysical
	case AttributeFrost:
		return AttributeIce
	default:
		return a
	}
}

// SVG returns the raw inline SVG markup string for the attribute.
func (a Attribute) SVG() string {
	return attributeSVGMap[a]
}

// IconURL returns the base64-encoded Data URI string containing the attribute's SVG icon.
func (a Attribute) IconURL() string {
	svg := a.SVG()
	if svg == "" {
		return ""
	}
	return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
}

// IconURL returns the official Enka CDN icon URL for the specialty.
func (s Specialty) IconURL() string {
	switch s {
	case SpecialtyAttack:
		return "https://enka.network/ui/zzz/IconAttack.png"
	case SpecialtyStun:
		return "https://enka.network/ui/zzz/IconStun.png"
	case SpecialtyAnomaly:
		return "https://enka.network/ui/zzz/IconAnomaly.png"
	case SpecialtySupport:
		return "https://enka.network/ui/zzz/IconSupport.png"
	case SpecialtyDefense:
		return "https://enka.network/ui/zzz/IconDefense.png"
	case SpecialtyRupture:
		return "https://enka.network/ui/zzz/IconRupture.png"
	default:
		return ""
	}
}

// IconURL returns the official Enka CDN icon URL for the rarity tier.
func (r Rarity) IconURL() string {
	switch r {
	case RarityS:
		return "https://enka.network/ui/zzz/ItemRarityS.png"
	case RarityA:
		return "https://enka.network/ui/zzz/ItemRarityA.png"
	case RarityB:
		return "https://enka.network/ui/zzz/ItemRarityB.png"
	default:
		return ""
	}
}

// Skin represents the equipped skin (outfit) of an agent.
type Skin struct {
	ID           int    `json:"id"`             // The internal ID of the skin.
	Name         string `json:"name"`           // The localized name of the skin.
	Description  string `json:"description"`    // The localized description of the skin.
	SplashArtURL string `json:"splash_art_url"` // The URL to the skin's splash art.
}

// SkillType represents a categorized category of combat skill.
type SkillType string

const (
	SkillTypeBasic   SkillType = "basic"   // Basic Attack
	SkillTypeDodge   SkillType = "dodge"   // Dodge & Counter
	SkillTypeAssist  SkillType = "assist"  // Quick & Defensive Assists
	SkillTypeSpecial SkillType = "special" // Special & EX Special Attack
	SkillTypeChain   SkillType = "chain"   // Chain Attack & Ultimate
	SkillTypePassive SkillType = "passive" // Core Passive & Additional Ability
)

// SkillParam represents a calculated numeric parameter or multiplier for a skill.
type SkillParam struct {
	Name  string `json:"name"`  // Localized parameter name.
	Value string `json:"value"` // Formatted value.
}

// Skill represents an agent's combat skill or passive ability.
type Skill struct {
	Level       int          `json:"level"`            // The level of the skill.
	Name        string       `json:"name"`             // The localized name of the skill.
	Description string       `json:"description"`      // The localized description of the skill.
	Type        SkillType    `json:"type"`             // Category type of the skill (basic, dodge, assist, special, chain, passive).
	TypeName    string       `json:"type_name"`        // Localized category type name.
	Params      []SkillParam `json:"params,omitempty"` // Numeric parameters / multiplier table.
}

// SkillGroup represents a categorized group of skills matching the in-game UI / Enka buttons.
type SkillGroup struct {
	Type     SkillType `json:"type"`      // Group category key ("basic", "special", "dodge", "chain", "assist", "passive").
	TypeName string    `json:"type_name"` // Localized category group name.
	Level    int       `json:"level"`     // Group level (1-12 for active skills, 0-6 for core passives).
	Skills   []Skill   `json:"skills"`    // Individual skills belonging to this category group.
}

// EvaluatedDescription returns the skill description with all scaling formulas ({CAL:...})
// evaluated for the skill's current level.
func (s Skill) EvaluatedDescription() string {
	return EvaluateFormulas(s.Description, s.Level)
}

// FormatHTML returns the skill description formatted as HTML with inline CSS colors,
// semantic icon spans, and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatHTML() string {
	return FormatHTML(s.Description, s.Level)
}

// FormatPlainText returns the skill description as clean plain text with all tags stripped
// and scaling formulas evaluated for the skill's current level.
func (s Skill) FormatPlainText() string {
	return FormatPlainText(s.Description, s.Level)
}

// FormatMarkdown returns the skill description formatted in Markdown (bold tags for colored values)
// with scaling formulas evaluated for the skill's current level.
func (s Skill) FormatMarkdown() string {
	return FormatMarkdown(s.Description, s.Level)
}

// Agent represents an enriched agent (character) showcased on a player's profile.
// A profile can showcase a maximum of 6 agents. It contains the agent's combat
// metadata, equipped gear, and final stats.
type Agent struct {
	ID                   int                 `json:"id"`                     // The internal ID of the agent.
	Name                 string              `json:"name"`                   // The localized name of the agent (e.g., "Ellen").
	Level                int                 `json:"level"`                  // The current level of the agent (1-60).
	Promotion            int                 `json:"promotion"`              // The promotion/ascension phase of the agent (0-5).
	MindscapeCinema      int                 `json:"mindscape_cinema"`       // The unlocked Mindscape Cinema level (0-6).
	CoreSkillEnhancement int                 `json:"core_skill_enhancement"` // The Core Skill enhancement level (0-6).
	Attribute            Attribute           `json:"attribute"`              // The elemental damage type (e.g., Ice).
	AttributeName        string              `json:"attribute_name"`         // The localized name of the attribute.
	Specialty            Specialty           `json:"specialty"`              // The combat role (e.g., Attack).
	SpecialtyName        string              `json:"specialty_name"`         // The localized name of the specialty.
	Rarity               Rarity              `json:"rarity"`                 // The rarity tier (S or A).
	Skin                 *Skin               `json:"skin"`                   // The currently equipped skin (can be nil if not found).
	SplashArtURL         string              `json:"splash_art_url"`         // The URL to the agent's splash art.
	Skills               []Skill             `json:"skills"`                 // The agent's skills and passives.
	Mindscapes           []MindscapeNode     `json:"mindscapes"`             // The agent's Mindscape Cinema levels (1-6).
	PotentialVision      *PotentialVision    `json:"potential_vision"`       // Potential Vision upgrade mechanics (can be nil if agent has none).
	WEngine              *WEngine            `json:"w_engine"`               // The currently equipped W-Engine (can be nil).
	DriveDiscs           []DriveDisc         `json:"drive_discs"`            // The equipped Drive Discs (up to 6).
	ActiveSetBonuses     []DriveDiscSetBonus `json:"active_set_bonuses"`     // The active 2-piece or 4-piece set bonuses.
	BaseStats            Stats               `json:"base_stats"`             // The agent's base combat stats before gear/buffs.
	Stats                Stats               `json:"stats"`                  // The agent's final combat stats including all gear/buffs.
}

// MindscapeNode represents a single Mindscape Cinema level (1-6) for an Agent.
type MindscapeNode struct {
	Rank        int    `json:"rank"`        // Cinema level (1 to 6).
	Name        string `json:"name"`        // Localized name of the Mindscape Cinema.
	Description string `json:"description"` // Localized description of the effect.
	Unlocked    bool   `json:"unlocked"`    // True if unlocked (MindscapeCinema >= Rank).
}

// FormatHTML returns the Mindscape description formatted as HTML with inline CSS colors.
func (m MindscapeNode) FormatHTML() string {
	return FormatHTML(m.Description)
}

// FormatPlainText returns the Mindscape description stripped of Rich Text formatting.
func (m MindscapeNode) FormatPlainText() string {
	return FormatPlainText(m.Description)
}

// FormatMarkdown returns the Mindscape description formatted with Markdown syntax.
func (m MindscapeNode) FormatMarkdown() string {
	return FormatMarkdown(m.Description)
}

// PotentialVision represents Potential Vision status and nodes for an Agent.
type PotentialVision struct {
	IsUnlocked bool                  `json:"is_unlocked"` // True if Potential Vision mechanic is unlocked.
	CurrentID  int                   `json:"current_id"`  // Current active Upgrade ID.
	Nodes      []PotentialVisionNode `json:"nodes"`       // All potential vision upgrade nodes.
}

// PotentialVisionNode represents a single Potential Vision upgrade node.
type PotentialVisionNode struct {
	ID          int    `json:"id"`          // Upgrade node ID.
	Level       int    `json:"level"`       // Level threshold (1 to 6).
	LevelName   string `json:"level_name"`  // Localized level title.
	Title       string `json:"title"`       // Localized title.
	Description string `json:"description"` // Localized effect description.
	IsActive    bool   `json:"is_active"`   // True if this node is active on the agent.
}

// FormatHTML returns the PotentialVisionNode description formatted as HTML with inline CSS colors.
func (p PotentialVisionNode) FormatHTML() string {
	return FormatHTML(p.Description)
}

// FormatPlainText returns the PotentialVisionNode description stripped of Rich Text formatting.
func (p PotentialVisionNode) FormatPlainText() string {
	return FormatPlainText(p.Description)
}

// FormatMarkdown returns the PotentialVisionNode description formatted with Markdown syntax.
func (p PotentialVisionNode) FormatMarkdown() string {
	return FormatMarkdown(p.Description)
}

// SubStatTotals calculates the sum of all sub-stats across all equipped Drive Discs.
// It groups them by PropertyID and sums the Rolls and Values.
// The returned slice is guaranteed to preserve the initial appearance order of sub-stats.
func (a *Agent) SubStatTotals() []StatValue {
	totals := make(map[PropertyID]StatValue)
	var order []PropertyID // Keep track of the order to ensure deterministic output

	for _, disc := range a.DriveDiscs {
		for _, sub := range disc.SubStats {
			if curr, exists := totals[sub.PropertyID]; exists {
				curr.Value += sub.Value
				curr.Rolls += sub.Rolls
				totals[sub.PropertyID] = curr
			} else {
				totals[sub.PropertyID] = sub
				order = append(order, sub.PropertyID)
			}
		}
	}

	result := make([]StatValue, 0, len(order))
	for _, id := range order {
		result = append(result, totals[id])
	}
	return result
}

// CountEffectiveRolls returns the total number of sub-stat rolls across all Drive Discs
// that match any of the provided target property IDs (also known as "effective" or "useful" rolls).
func (a *Agent) CountEffectiveRolls(targetProps ...PropertyID) int {
	total := 0
	targetMap := make(map[PropertyID]bool)
	for _, p := range targetProps {
		targetMap[p] = true
	}

	for _, disc := range a.DriveDiscs {
		for _, sub := range disc.SubStats {
			if targetMap[sub.PropertyID] {
				total += sub.Rolls
			}
		}
	}
	return total
}

func getStatName(st store.MetadataStore, key string, lang Language) string {
	if st != nil {
		if val := st.Localize(key, string(lang)); val != "" && val != key {
			return val
		}
	}
	return key
}

// FormattedUIStats generates a complete breakdown of base vs added stats for UI display.
// Stat names are localized according to the optional lang parameter (defaults to LangEN).
// This structure precisely matches the visual representation and layout seen in the in-game
// stat panel or on platforms like Enka.Network.
func (a *Agent) FormattedUIStats(lang ...Language) UIStats {
	l := LangEN
	if len(lang) > 0 {
		l = lang[0]
	}

	st, _ := store.Default()

	hpName := getStatName(st, locKeyHP, l)
	atkName := getStatName(st, locKeyATK, l)
	defName := getStatName(st, locKeyDEF, l)
	impactName := getStatName(st, locKeyImpact, l)
	critRateName := getStatName(st, locKeyCritRate, l)
	critDMGName := getStatName(st, locKeyCritDMG, l)
	anomalyMasteryName := getStatName(st, locKeyAnomalyMastery, l)
	anomalyProficiencyName := getStatName(st, locKeyAnomalyProficiency, l)
	penRatioName := getStatName(st, locKeyPenRatio, l)
	penFlatName := getStatName(st, locKeyPenFlat, l)
	energyRegenName := getStatName(st, locKeyEnergyRegen, l)
	sheerForceName := getStatName(st, locKeySheerForce, l)

	var attrDMGName string
	var attrDMGProp PropertyID
	switch a.Attribute.BaseAttribute() {
	case AttributePhysical:
		attrDMGName = getStatName(st, locKeyPhysicalDMGBonus, l)
		attrDMGProp = PropPhysicalDMGBonus
	case AttributeFire:
		attrDMGName = getStatName(st, locKeyFireDMGBonus, l)
		attrDMGProp = PropFireDMGBonus
	case AttributeIce:
		attrDMGName = getStatName(st, locKeyIceDMGBonus, l)
		attrDMGProp = PropIceDMGBonus
	case AttributeElectric:
		attrDMGName = getStatName(st, locKeyElectricDMGBonus, l)
		attrDMGProp = PropElectricDMGBonus
	case AttributeEther:
		attrDMGName = getStatName(st, locKeyEtherDMGBonus, l)
		attrDMGProp = PropEtherDMGBonus
	case AttributeWind:
		attrDMGName = getStatName(st, locKeyWindDMGBonus, l)
		attrDMGProp = PropWindDMGBonus
	}

	return UIStats{
		HP:                 formatBreakdown(PropBaseHP, hpName, a.BaseStats.HP, a.Stats.HP, false, 0, 1),
		ATK:                formatBreakdown(PropBaseATK, atkName, a.BaseStats.ATK, a.Stats.ATK, false, 0, 1),
		DEF:                formatBreakdown(PropBaseDEF, defName, a.BaseStats.DEF, a.Stats.DEF, false, 0, 1),
		Impact:             formatBreakdown(PropBaseImpact, impactName, a.BaseStats.Impact, a.Stats.Impact, false, 0, 1),
		CritRate:           formatBreakdown(PropBaseCritRate, critRateName, a.BaseStats.CritRate, a.Stats.CritRate, true, 0, 1),
		CritDMG:            formatBreakdown(PropBaseCritDMG, critDMGName, a.BaseStats.CritDMG, a.Stats.CritDMG, true, 0, 1),
		AttributeDMGBonus:  formatBreakdown(attrDMGProp, attrDMGName, a.BaseStats.AttributeDMGBonus, a.Stats.AttributeDMGBonus, true, 0, 1),
		AnomalyMastery:     formatBreakdown(PropBaseAnomalyMastery, anomalyMasteryName, a.BaseStats.AnomalyMastery, a.Stats.AnomalyMastery, false, 0, 1),
		AnomalyProficiency: formatBreakdown(PropBaseAnomalyProficiency, anomalyProficiencyName, a.BaseStats.AnomalyProficiency, a.Stats.AnomalyProficiency, false, 0, 1),
		PenRatio:           formatBreakdown(PropBasePENRatio, penRatioName, a.BaseStats.PenRatio, a.Stats.PenRatio, true, 0, 1),
		PenFlat:            formatBreakdown(PropBasePENFlat, penFlatName, a.BaseStats.PenFlat, a.Stats.PenFlat, false, 0, 1),
		EnergyRegen:        formatBreakdown(PropBaseEnergyRegen, energyRegenName, a.BaseStats.EnergyRegen, a.Stats.EnergyRegen, false, 2, 2),
		SheerForce:         formatBreakdown(PropBaseSheerForce, sheerForceName, a.BaseStats.SheerForce, a.Stats.SheerForce, false, 0, 1),
	}
}

// GroupedSkills returns the agent's skills categorized into 6 distinct groups
// matching the in-game skill buttons:
// 1. Passives / Talents (Core Passive + Additional Ability)
// 2. Basic Attack
// 3. Dodge
// 4. Assist
// 5. Special Attack
// 6. Chain Attack & Ultimate
func (a *Agent) GroupedSkills() []SkillGroup {
	order := []SkillType{
		SkillTypePassive,
		SkillTypeBasic,
		SkillTypeDodge,
		SkillTypeChain,
		SkillTypeAssist,
		SkillTypeSpecial,
	}

	groupsMap := make(map[SkillType]*SkillGroup)
	for _, sk := range a.Skills {
		st := sk.Type
		if st == "" {
			st = SkillTypeBasic
		}
		grp, ok := groupsMap[st]
		if !ok {
			grp = &SkillGroup{
				Type:     st,
				TypeName: sk.TypeName,
				Level:    sk.Level,
				Skills:   make([]Skill, 0),
			}
			groupsMap[st] = grp
		}
		grp.Skills = append(grp.Skills, sk)
	}

	result := make([]SkillGroup, 0, len(order))
	for _, st := range order {
		if grp, ok := groupsMap[st]; ok {
			result = append(result, *grp)
		}
	}
	return result
}
