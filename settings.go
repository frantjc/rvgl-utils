package rvglutils

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	xslices "github.com/frantjc/x/slices"
	"github.com/pelletier/go-toml/v2"
)

type ResolveSettingsINIOpts struct {
	Profile  string
	PathList string
}

func (o *ResolveSettingsINIOpts) Apply(opts *ResolveSettingsINIOpts) {
	if o != nil {
		if opts != nil {
			if o.Profile != "" {
				opts.Profile = o.Profile
			}
			if o.PathList != "" {
				opts.PathList = o.PathList
			}
		}
	}
}

type ResolveSettingsINIOpt interface {
	Apply(*ResolveSettingsINIOpts)
}

type Settings struct {
	Video   VideoSettings   `toml:"Video" ini:"Video"`
	Audio   AudioSettings   `toml:"Audio" ini:"Audio"`
	Network NetworkSettings `toml:"Network" ini:"Network"`
	Editor  EditorSettings  `toml:"Editor" ini:"Editor"`
	Misc    MiscSettings    `toml:"Misc" ini:"Misc"`
}

type VideoSettings struct {
	EnvLevel         int `toml:"EnvLevel" ini:"EnvLevel"`
	ParticleLevel    int `toml:"ParticleLevel" ini:"ParticleLevel"`
	MirrorFlag       int `toml:"MirrorFlag" ini:"MirrorFlag"`
	ShadowFlag       int `toml:"ShadowFlag" ini:"ShadowFlag"`
	LightFlag        int `toml:"LightFlag" ini:"LightFlag"`
	RippleFlag       int `toml:"RippleFlag" ini:"RippleFlag"`
	InstanceFlag     int `toml:"InstanceFlag" ini:"InstanceFlag"`
	SkidFlag         int `toml:"SkidFlag" ini:"SkidFlag"`
	EffectFlag       int `toml:"EffectFlag" ini:"EffectFlag"`
	ScreenWidth      int `toml:"ScreenWidth" ini:"ScreenWidth"`
	ScreenHeight     int `toml:"ScreenHeight" ini:"ScreenHeight"`
	ScreenBpp        int `toml:"ScreenBpp" ini:"ScreenBpp"`
	ScreenHz         int `toml:"ScreenHz" ini:"ScreenHz"`
	Orientation      int `toml:"Orientation" ini:"Orientation"`
	Brightness       int `toml:"Brightness" ini:"Brightness"`
	Contrast         int `toml:"Contrast" ini:"Contrast"`
	DrawDist         int `toml:"DrawDist" ini:"DrawDist"`
	Texture32        int `toml:"Texture32" ini:"Texture32"`
	Vsync            int `toml:"Vsync" ini:"Vsync"`
	ShowFPS          int `toml:"ShowFPS" ini:"ShowFPS"`
	LimitFPS         int `toml:"LimitFPS" ini:"LimitFPS"`
	LimitLatency     int `toml:"LimitLatency" ini:"LimitLatency"`
	CenterHud        int `toml:"CenterHud" ini:"CenterHud"`
	SplitScreenMode  int `toml:"SplitScreenMode" ini:"SplitScreenMode"`
	MaintainAspect   int `toml:"MaintainAspect" ini:"MaintainAspect"`
	TextureFilter    int `toml:"TextureFilter" ini:"TextureFilter"`
	MipmapFilter     int `toml:"MipmapFilter" ini:"MipmapFilter"`
	Anisotropy       int `toml:"Anisotropy" ini:"Anisotropy"`
	Antialias        int `toml:"Antialias" ini:"Antialias"`
	GenerateMipmaps  int `toml:"GenerateMipmaps" ini:"GenerateMipmaps"`
	CompressTextures int `toml:"CompressTextures" ini:"CompressTextures"`
	Threaded         int `toml:"Threaded" ini:"Threaded"`
	Profile          int `toml:"Profile" ini:"Profile"`
	Shaders          int `toml:"Shaders" ini:"Shaders"`
	ShaderLights     int `toml:"ShaderLights" ini:"ShaderLights"`
	ShaderShadows    int `toml:"ShaderShadows" ini:"ShaderShadows"`
	ShaderEffects    int `toml:"ShaderEffects" ini:"ShaderEffects"`
	EnableFBO        int `toml:"EnableFBO" ini:"EnableFBO"`
	EnableUBO        int `toml:"EnableUBO" ini:"EnableUBO"`
	EnableSSO        int `toml:"EnableSSO" ini:"EnableSSO"`
	SortLevel        int `toml:"SortLevel" ini:"SortLevel"`
	Compositor       int `toml:"Compositor" ini:"Compositor"`
	HighDPI          int `toml:"HighDPI" ini:"HighDPI"`
}

type AudioSettings struct {
	SfxVol      int `toml:"SfxVol" ini:"SfxVol"`
	MusicVol    int `toml:"MusicVol" ini:"MusicVol"`
	SfxChannels int `toml:"SfxChannels" ini:"SfxChannels"`
	SampleRate  int `toml:"SampleRate" ini:"SampleRate"`
	MusicOn     int `toml:"MusicOn" ini:"MusicOn"`
}

type NetworkSettings struct {
	HostComputer string `toml:"HostComputer" ini:"HostComputer"`
	LocalPort    int    `toml:"LocalPort" ini:"LocalPort"`
	DiscordRPC   int    `toml:"DiscordRPC" ini:"DiscordRPC"`
}

type EditorSettings struct {
	LegacyCameraControls int `toml:"LegacyCameraControls" ini:"LegacyCameraControls"`
	LegacyEditorControls int `toml:"LegacyEditorControls" ini:"LegacyEditorControls"`
}

type MiscSettings struct {
	UseProfiles    int    `toml:"UseProfiles" ini:"UseProfiles"`
	DefaultProfile string `toml:"DefaultProfile" ini:"DefaultProfile"`
	Frontend       string `toml:"Frontend" ini:"Frontend"`
	Version        string `toml:"Version" ini:"Version"`
}

type ProfileSettings struct {
	Game        ProfileGameSettings       `toml:"Game" ini:"Game"`
	Controller1 ProfileControllerSettings `toml:"Controller1" ini:"Controller1"`
	Controller2 ProfileControllerSettings `toml:"Controller2" ini:"Controller2"`
	Controller3 ProfileControllerSettings `toml:"Controller3" ini:"Controller3"`
	Controller4 ProfileControllerSettings `toml:"Controller4" ini:"Controller4"`
}

type ProfileGameSettings struct {
	Pickups         int    `toml:"Pickups" ini:"Pickups"`
	NCars           int    `toml:"NCars" ini:"NCars"`
	NLaps           int    `toml:"NLaps" ini:"NLaps"`
	PlayMode        int    `toml:"PlayMode" ini:"PlayMode"`
	Language        string `toml:"Language" ini:"Language"`
	LevelDir        string `toml:"LevelDir" ini:"LevelDir"`
	PlayerName1     string `toml:"PlayerName1" ini:"PlayerName1"`
	CarDir1         string `toml:"CarDir1" ini:"CarDir1"`
	PlayerName2     string `toml:"PlayerName2" ini:"PlayerName2"`
	CarDir2         string `toml:"CarDir2" ini:"CarDir2"`
	PlayerName3     string `toml:"PlayerName3" ini:"PlayerName3"`
	CarDir3         string `toml:"CarDir3" ini:"CarDir3"`
	PlayerName4     string `toml:"PlayerName4" ini:"PlayerName4"`
	CarDir4         string `toml:"CarDir4" ini:"CarDir4"`
	SpeedUnits      int    `toml:"SpeedUnits" ini:"SpeedUnits"`
	LocalCamera     int    `toml:"LocalCamera" ini:"LocalCamera"`
	DrawRearView    int    `toml:"DrawRearView" ini:"DrawRearView"`
	AnimateRearView int    `toml:"AnimateRearView" ini:"AnimateRearView"`
	RearViewType    int    `toml:"RearViewType" ini:"RearViewType"`
	GhostType       int    `toml:"GhostType" ini:"GhostType"`
	ShowGhost       int    `toml:"ShowGhost" ini:"ShowGhost"`
	CatchUp         int    `toml:"CatchUp" ini:"CatchUp"`
	CupDC           int    `toml:"CupDC" ini:"CupDC"`
	Difficulty      int    `toml:"Difficulty" ini:"Difficulty"`
	BattleTagTime   int    `toml:"BattleTagTime" ini:"BattleTagTime"`
	NumberOfPlayers int    `toml:"NumberOfPlayers" ini:"NumberOfPlayers"`
	MultiPlayerCPU  int    `toml:"MultiPlayerCPU" ini:"MultiPlayerCPU"`
	FinalLapMsg     int    `toml:"FinalLapMsg" ini:"FinalLapMsg"`
	FinalCam        int    `toml:"FinalCam" ini:"FinalCam"`
	WeaponCam       int    `toml:"WeaponCam" ini:"WeaponCam"`
	RandomSkins     int    `toml:"RandomSkins" ini:"RandomSkins"`
}

type ProfileControllerSettings struct {
	Joystick          int    `toml:"Joystick" ini:"Joystick"`
	ForceFeedback     int    `toml:"ForceFeedback" ini:"ForceFeedback"`
	NonLinearSteering int    `toml:"NonLinearSteering" ini:"NonLinearSteering"`
	SteeringDeadzone  int    `toml:"SteeringDeadzone" ini:"SteeringDeadzone"`
	SteeringRange     int    `toml:"SteeringRange" ini:"SteeringRange"`
	ButtonOpacity     int    `toml:"ButtonOpacity" ini:"ButtonOpacity"`
	KeyLeft           uint32 `toml:"KeyLeft" ini:"KeyLeft"`
	KeyRight          uint32 `toml:"KeyRight" ini:"KeyRight"`
	KeyFwd            uint32 `toml:"KeyFwd" ini:"KeyFwd"`
	KeyBack           uint32 `toml:"KeyBack" ini:"KeyBack"`
	KeyFire           uint32 `toml:"KeyFire" ini:"KeyFire"`
	KeyReset          uint32 `toml:"KeyReset" ini:"KeyReset"`
	KeyReposition     uint32 `toml:"KeyReposition" ini:"KeyReposition"`
	KeyHonka          uint32 `toml:"KeyHonka" ini:"KeyHonka"`
	KeyChangeCamera   uint32 `toml:"KeyChangeCamera" ini:"KeyChangeCamera"`
	KeyRearView       uint32 `toml:"KeyRearView" ini:"KeyRearView"`
	KeyPause          uint32 `toml:"KeyPause" ini:"KeyPause"`
}

func newLoadConfigOpts(opts ...ResolveSettingsINIOpt) *ResolveSettingsINIOpts {
	o := &ResolveSettingsINIOpts{
		PathList: strings.Join(xslices.Map(strings.Split(DefaultPrefPathList, string(os.PathListSeparator)), func(prefPath string, _ int) string {
			return filepath.Join(prefPath, "profiles")
		}), string(os.PathListSeparator)),
	}

	for _, opt := range opts {
		opt.Apply(o)
	}

	return o
}

func ResolveSettingsINI(opts ...ResolveSettingsINIOpt) (string, error) {
	var (
		o    = newLoadConfigOpts(opts...)
		dirs = strings.Split(o.PathList, string(os.PathListSeparator))
	)

	cwd, err := os.Getwd()
	if err == nil {
		dirs = append(dirs, cwd)
	}

	for _, dir := range dirs {
		name := filepath.Join(dir, "rvgl.ini")
		if o.Profile != "" {
			name = filepath.Join(dir, strings.ToLower(o.Profile), "profile.ini")
		}

		if _, err := os.Stat(name); err != nil {
			continue
		}

		return name, nil
	}

	return "", fmt.Errorf("load settings .ini")
}

type skipLinesReader struct {
	scanner *bufio.Scanner
	prefix  []byte
	buf     *bytes.Buffer
}

func (r *skipLinesReader) Read(p []byte) (int, error) {
	if r.buf == nil {
		r.buf = new(bytes.Buffer)
	}

	var (
		lenBuf = r.buf.Len()
		lenP   = len(p)
	)
	if lenBuf < lenP {
		for r.scanner.Scan() {
			if bytes.HasPrefix(r.scanner.Bytes(), r.prefix) {
				continue
			}

			n, err := r.buf.Write(append(r.scanner.Bytes(), '\n'))
			if err != nil {
				return 0, err
			} else if n+lenBuf > lenP {
				break
			}
		}
	}

	return r.buf.Read(p)
}

func DecodeSettingsINI(r io.Reader) (*Settings, error) {
	s := &Settings{}
	return s, toml.NewDecoder(&skipLinesReader{scanner: bufio.NewScanner(r), prefix: []byte(";")}).Decode(s)
}

func DecodeProfileSettingsINI(r io.Reader) (*ProfileSettings, error) {
	s := &ProfileSettings{}
	return s, toml.NewDecoder(&skipLinesReader{scanner: bufio.NewScanner(r), prefix: []byte(";")}).Decode(s)
}

var (
	//go:embed settings-ini-header.txt
	settingsINIHeader []byte
)

func EncodeSettingsINI(w io.Writer, v any) error {
	if _, err := w.Write(settingsINIHeader); err != nil {
		return err
	}

	return toml.NewEncoder(w).Encode(v)
}
