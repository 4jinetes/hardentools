// Hardentools
// Copyright (C) 2017  Security Without Borders
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

// available office versions
var standardOfficeVersions = []string{
	"12.0", // Office 2007
	"14.0", // Office 2010
	"15.0", // Office 2013
	"16.0", // Office 2016
}

// standard office apps that are hardened
var standardOfficeApps = []string{"Excel", "PowerPoint", "Word"}

// OfficeRegistryRegExSingleDWORD is the data type for a RegEx Path / Single Value DWORD combination
type OfficeRegistryRegExSingleDWORD struct {
	RootKey        registry.Key
	PathRegEx      string
	ValueName      string
	HardenedValue  uint32
	OfficeApps     []string
	OfficeVersions []string
	shortName      string
	longName       string
	description    string
}

// OfficeOLE hardens Office Packager Objects
// 0 - No prompt from Office when user clicks, object executes
// 1 - Prompt from Office when user clicks, object executes
// 2 - No prompt, Object does not execute
var OfficeOLE = &OfficeRegistryRegExSingleDWORD{
	RootKey:        registry.CURRENT_USER,
	PathRegEx:      "SOFTWARE\\Microsoft\\Office\\%s\\%s\\Security",
	ValueName:      "PackagerPrompt",
	HardenedValue:  2,
	OfficeApps:     standardOfficeApps,
	OfficeVersions: standardOfficeVersions,
	shortName:      "OfficeOLE",
	longName:       "Office Packager Objects (OLE)"}

// OfficeMacros contains Macro registry keys
// 1 - Enable all
// 2 - Disable with notification
// 3 - Digitally signed only
// 4 - Disable all
var OfficeMacros = &OfficeRegistryRegExSingleDWORD{
	RootKey:        registry.CURRENT_USER,
	PathRegEx:      "SOFTWARE\\Microsoft\\Office\\%s\\%s\\Security",
	ValueName:      "VBAWarnings",
	HardenedValue:  4,
	OfficeApps:     standardOfficeApps,
	OfficeVersions: standardOfficeVersions,
	shortName:      "OfficeMacros",
	longName:       "Office Macros"}

// OfficeActiveX contains ActiveX registry keys
var OfficeActiveX = &RegistrySingleValueDWORD{
	RootKey:       registry.CURRENT_USER,
	Path:          "SOFTWARE\\Microsoft\\Office\\Common\\Security",
	ValueName:     "DisableAllActiveX",
	HardenedValue: 1,
	shortName:     "OfficeActiveX",
	longName:      "Office ActiveX"}

//// DDE Mitigations for Word, Outlook and Excel
// Doesn't harden OneNote for now (due to high impact).
// [HKEY_CURRENT_USER\Software\Microsoft\Office\%s\Word\Options]
// [HKEY_CURRENT_USER\Software\Microsoft\Office\%s\Word\Options\WordMail] (this one is for Outlook)
// [HKEY_CURRENT_USER\Software\Microsoft\Office\%s\Excel\Options]
//    "DontUpdateLinks"=dword:00000001
//
// additionally only for Excel:
// [HKEY_CURRENT_USER\Software\Microsoft\Office\%s\Excel\Options]
//   "DDEAllowed"=dword:00000000
//   "DDECleaned"=dword:00000001
//   "Options"=dword:00000117
// [HKEY_CURRENT_USER\Software\Microsoft\Office\<version>\Excel\Security]
//   WorkbookLinkWarnings(DWORD) = 2
//
// for Word&Outlook 2007:
// [HKEY_CURRENT_USER\Software\Microsoft\Office\12.0\Word\Options\vpref]
//    fNoCalclinksOnopen_90_1(DWORD)=1
var pathRegExOptions = "SOFTWARE\\Microsoft\\Office\\%s\\%s\\Options"
var pathRegExWordMail = "SOFTWARE\\Microsoft\\Office\\%s\\%s\\Options\\WordMail"
var pathRegExSecurity = "Software\\Microsoft\\Office\\%s\\%s\\Security"
var pathWord2007 = "Software\\Microsoft\\Office\\12.0\\Word\\Options\\vpref"

// OfficeDDE contains the registry keys for DDE hardening
var OfficeDDE = &MultiHardenInterfaces{
	hardenInterfaces: []HardenInterface{
		&OfficeRegistryRegExSingleDWORD{
			RootKey:       registry.CURRENT_USER,
			PathRegEx:     pathRegExOptions,
			ValueName:     "DontUpdateLinks",
			HardenedValue: 1,
			OfficeApps:    []string{"Word", "Excel"},
			OfficeVersions: []string{
				"14.0", // Office 2010
				"15.0", // Office 2013
				"16.0", // Office 2016
			},
			shortName: "OfficeDDE_DontUpdateLinksWordExcel"},

		&OfficeRegistryRegExSingleDWORD{
			RootKey:       registry.CURRENT_USER,
			PathRegEx:     pathRegExWordMail,
			ValueName:     "DontUpdateLinks",
			HardenedValue: 1,
			OfficeApps:    []string{"Word"},
			OfficeVersions: []string{
				"14.0", // Office 2010
				"15.0", // Office 2013
				"16.0", // Office 2016
			},
			shortName: "OfficeDDE_DontUpdateLinksWordMail"},

		&OfficeRegistryRegExSingleDWORD{
			RootKey:        registry.CURRENT_USER,
			PathRegEx:      pathRegExOptions,
			ValueName:      "DDEAllowed",
			HardenedValue:  0,
			OfficeApps:     []string{"Excel"},
			OfficeVersions: standardOfficeVersions,
			shortName:      "OfficeDDE_DDEAllowedExcel"},

		&OfficeRegistryRegExSingleDWORD{
			RootKey:        registry.CURRENT_USER,
			PathRegEx:      pathRegExOptions,
			ValueName:      "DDECleaned",
			HardenedValue:  1,
			OfficeApps:     []string{"Excel"},
			OfficeVersions: standardOfficeVersions,
			shortName:      "OfficeDDE_DDECleanedExcel"},

		&OfficeRegistryRegExSingleDWORD{
			RootKey:        registry.CURRENT_USER,
			PathRegEx:      pathRegExOptions,
			ValueName:      "Options",
			HardenedValue:  0x117,
			OfficeApps:     []string{"Excel"},
			OfficeVersions: standardOfficeVersions,
			shortName:      "OfficeDDE_OptionsExcel"},

		&OfficeRegistryRegExSingleDWORD{
			RootKey:        registry.CURRENT_USER,
			PathRegEx:      pathRegExSecurity,
			ValueName:      "WorkbookLinkWarnings",
			HardenedValue:  2,
			OfficeApps:     []string{"Excel"},
			OfficeVersions: standardOfficeVersions,
			shortName:      "OfficeDDE_WorkbookLinksExcel"},

		&RegistrySingleValueDWORD{
			RootKey:       registry.CURRENT_USER,
			Path:          pathWord2007,
			ValueName:     "fNoCalclinksOnopen_90_1",
			HardenedValue: 1,
			shortName:     "OfficeDDE_Word2007"},
	},
	shortName: "OfficeDDE",
	longName:  "Office DDE  Links",
}

//// HardenInterface methods

// Harden hardens OfficeRegistryRegExSingleDWORD registry values
func (officeRegEx OfficeRegistryRegExSingleDWORD) Harden(harden bool) error {

	for _, officeVersion := range officeRegEx.OfficeVersions {
		for _, officeApp := range officeRegEx.OfficeApps {
			path := fmt.Sprintf(officeRegEx.PathRegEx, officeVersion, officeApp)

			// build a RegistrySingleValueDWORD so we can reuse the Harden() method
			var singleDWORD = &RegistrySingleValueDWORD{
				RootKey:       officeRegEx.RootKey,
				Path:          path,
				ValueName:     officeRegEx.ValueName,
				HardenedValue: officeRegEx.HardenedValue,
				shortName:     officeRegEx.shortName,
				longName:      officeRegEx.longName,
				description:   officeRegEx.description,
			}

			// call RegistrySingleValueDWORD Harden method to Harden or Restore.
			err := singleDWORD.Harden(harden)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// IsHardened verifies if OfficeRegistryRegExSingleDWORD is already hardenend
func (officeRegEx OfficeRegistryRegExSingleDWORD) IsHardened() bool {
	var hardened = true

	for _, officeVersion := range officeRegEx.OfficeVersions {
		for _, officeApp := range officeRegEx.OfficeApps {
			path := fmt.Sprintf(officeRegEx.PathRegEx, officeVersion, officeApp)

			// build a RegistrySingleValueDWORD so we can reuse the isHardened() method
			var singleDWORD = &RegistrySingleValueDWORD{
				RootKey:       officeRegEx.RootKey,
				Path:          path,
				ValueName:     officeRegEx.ValueName,
				HardenedValue: officeRegEx.HardenedValue,
			}

			if !singleDWORD.IsHardened() {
				hardened = false
			}
		}
	}
	return hardened
}

// Name returns the (short) name of the harden item
func (officeRegEx OfficeRegistryRegExSingleDWORD) Name() string {
	return officeRegEx.shortName
}

// LongName returns the long name of the harden item
func (officeRegEx OfficeRegistryRegExSingleDWORD) LongName() string {
	return officeRegEx.longName
}

// Description of the harden item
func (officeRegEx OfficeRegistryRegExSingleDWORD) Description() string {
	return officeRegEx.description
}
