package commUtils

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	commonUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	fileUtils "github.com/aaronchen2k/deeptest/pkg/lib/file"
	"os"
	"path/filepath"
	"strings"
)

func GetExecDir() (dir string) { // where exe file in
	exeDir, _ := os.Executable()

	if commonUtils.IsRelease() { // release
		dir = filepath.Dir(exeDir)
	} else { // debug mode
		if strings.Index(strings.ToLower(exeDir), "goland") > -1 { // run with ide
			dir = os.Getenv("ZTF_CODE_DIR")
		} else {
			dir = GetWorkDir()
		}
	}

	dir, _ = filepath.Abs(dir)
	dir = fileUtils.AddSepIfNeeded(dir)

	return
}

func GetWorkDir() (dir string) {
	home, _ := fileUtils.GetUserHome()
	dir = filepath.Join(home, consts.App)
	dir = fileUtils.AddSepIfNeeded(dir)
	fileUtils.MkDirIfNeeded(dir)

	return
}
