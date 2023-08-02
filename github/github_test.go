package github

import "testing"

func TestGetReleasesInfo(t *testing.T) {
	tInfo, err := GetReleasesInfo("git-for-windows/git")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(tInfo.Releases, tInfo.Tag)
		//f := tInfo.Releases.Array()[0]
	}
}

func TestGetGetReleasesEx(t *testing.T) {
	turl, err := GetReleasesEx("git-for-windows/git", ".*32-bit.exe")
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(turl)
	}
}
