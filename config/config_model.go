package config

type AppConfigYAML struct {
	Port  int    `yaml:"port"`
	Salt  string `yaml:"salt"`
	Image struct {
		StoragePath    string `yaml:"storage_path"`
		ResponsePath   string `yaml:"response_path"`
		MaxWidth       int    `yaml:"maxwidth"`
		MaxHeight      int    `yaml:"maxheight"`
		Quality        int    `yaml:"quality"`
		PostProcessing struct {
			Sharpen    float64 `yaml:"sharpen"`
			Brightness float64 `yaml:"brightness"`
			Contrast   float64 `yaml:"contrast"`
			Gamma      float64 `yaml:"gamma"`
		} `yaml:"post_processing"`
		Downscale string `yaml:"downscale"`
	}
	Thumbnail struct {
		StoragePath    string `yaml:"storage_path"`
		ResponsePath   string `yaml:"response_path"`
		MaxWidth       int    `yaml:"maxwidth"`
		MaxHeight      int    `yaml:"maxheight"`
		Quality        int    `yaml:"quality"`
		PostProcessing struct {
			Sharpen    float64 `yaml:"sharpen"`
			Brightness float64 `yaml:"brightness"`
			Contrast   float64 `yaml:"contrast"`
			Gamma      float64 `yaml:"gamma"`
		} `yaml:"post_processing"`
		Resize struct {
			Upscale   string `yaml:"upscale"`
			Downscale string `yaml:"downscale"`
		} `yaml:"resize"`
	}
	Logging struct {
		Path string `yaml:"path"`
	}
}
