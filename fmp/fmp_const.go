package fmp

type FmpError string

func (e FmpError) Error() string { return string(e) }

var (
	ErrRead               = FmpError("read error")
	ErrBadMagic           = FmpError("bad magic number")
	ErrBadHeader          = FmpError("bad header")
	ErrUnsupportedCharset = FmpError("unsupported character set")
	ErrBadSectorCount     = FmpError("bad sector count")
	ErrBadSectorHeader    = FmpError("bad sector header")
	ErrBadChunk           = FmpError("bad chunk")
)

type FmpChunkType uint8

const (
	FmpChunkSimpleData FmpChunkType = iota
	FmpChunkSegmentedData
	FmpChunkSimpleKeyValue
	FmpChunkLongKeyValue
	FmpChunkPathPush
	FmpChunkPathPushLong
	FmpChunkPathPop
	FmpChunkNoop
)

type FmpFieldType uint8

const (
	FmpFieldSimple      FmpFieldType = 1
	FmpFieldCalculation FmpFieldType = 2
	FmpFieldScript      FmpFieldType = 3
)

type FmpFieldStorageType uint8

const (
	FmpFieldStorageRegular             FmpFieldStorageType = 0
	FmpFieldStorageGlobal              FmpFieldStorageType = 1
	FmpFieldStorageCalculation         FmpFieldStorageType = 8
	FmpFieldStorageUnstoredCalculation FmpFieldStorageType = 10
)

type FmpDataType uint8

const (
	FmpDataText      FmpDataType = 1
	FmpDataNumber    FmpDataType = 2
	FmpDataDate      FmpDataType = 3
	FmpDataTime      FmpDataType = 4
	FmpDataTS        FmpDataType = 5
	FmpDataContainer FmpDataType = 6
)

type FmpAutoEnterOption uint8

const (
	FmpAutoEnterData FmpAutoEnterOption = iota
	FmpAutoEnterSerialNumber
	FmpAutoEnterCalculation
	FmpAutoEnterCalculationReplacingExistingValue
	FmpAutoEnterFromLastVisitedRecord
	FmpAutoEnterCreateDate
	FmpAutoEnterCreateTime
	FmpAutoEnterCreateTS
	FmpAutoEnterCreateName
	FmpAutoEnterCreateAccountName
	FmpAutoEnterModDate
	FmpAutoEnterModTime
	FmpAutoEnterModTS
	FmpAutoEnterModName
	FmpAutoEnterModAccountName
)

var autoEnterPresetMap = map[uint8]FmpAutoEnterOption{
	0: FmpAutoEnterCreateDate,
	1: FmpAutoEnterCreateTime,
	2: FmpAutoEnterCreateTS,
	3: FmpAutoEnterCreateName,
	4: FmpAutoEnterCreateAccountName,
	5: FmpAutoEnterModDate,
	6: FmpAutoEnterModTime,
	7: FmpAutoEnterModTS,
	8: FmpAutoEnterModName,
	9: FmpAutoEnterModAccountName,
}

var autoEnterOptionMap = map[uint8]FmpAutoEnterOption{
	2:   FmpAutoEnterSerialNumber,
	4:   FmpAutoEnterData,
	8:   FmpAutoEnterCalculation,
	16:  FmpAutoEnterFromLastVisitedRecord,
	32:  FmpAutoEnterCalculation,
	136: FmpAutoEnterCalculationReplacingExistingValue,
}

type FmpScriptStepType uint64

const (
	FmpScriptPerformScript                     FmpScriptStepType = 1
	FmpScriptSaveCopyAsXML                     FmpScriptStepType = 3
	FmpScriptGoToNextField                     FmpScriptStepType = 4
	FmpScriptGoToPreviousField                 FmpScriptStepType = 5
	FmpScriptGoToLayout                        FmpScriptStepType = 6
	FmpScriptNewRecordRequest                  FmpScriptStepType = 7
	FmpScriptDuplicateRecordRequest            FmpScriptStepType = 8
	FmpScriptDeleteRecordRequest               FmpScriptStepType = 9
	FmpScriptDeleteAllRecords                  FmpScriptStepType = 10
	FmpScriptInsertFromIndex                   FmpScriptStepType = 11
	FmpScriptInsertFromLastVisited             FmpScriptStepType = 12
	FmpScriptInsertCurrentDate                 FmpScriptStepType = 13
	FmpScriptInsertCurrentTime                 FmpScriptStepType = 14
	FmpScriptGoToRecordRequestPage             FmpScriptStepType = 16
	FmpScriptGoToField                         FmpScriptStepType = 17
	FmpScriptCheckSelection                    FmpScriptStepType = 18
	FmpScriptCheckRecord                       FmpScriptStepType = 19
	FmpScriptCheckFoundSet                     FmpScriptStepType = 20
	FmpScriptUnsortRecords                     FmpScriptStepType = 21
	FmpScriptEnterFindMode                     FmpScriptStepType = 22
	FmpScriptShowAllRecords                    FmpScriptStepType = 23
	FmpScriptModifyLastFind                    FmpScriptStepType = 24
	FmpScriptOmitRecord                        FmpScriptStepType = 25
	FmpScriptOmitMultipleRecords               FmpScriptStepType = 26
	FmpScriptShowOmmitedOnly                   FmpScriptStepType = 27
	FmpScriptPerformFind                       FmpScriptStepType = 28
	FmpScriptShowHideToolbars                  FmpScriptStepType = 29
	FmpScriptViewAs                            FmpScriptStepType = 30
	FmpScriptAdjustWindow                      FmpScriptStepType = 31
	FmpScriptOpenHelp                          FmpScriptStepType = 32
	FmpScriptOpenFile                          FmpScriptStepType = 33
	FmpScriptCloseFile                         FmpScriptStepType = 34
	FmpScriptImportRecords                     FmpScriptStepType = 35
	FmpScriptExportRecords                     FmpScriptStepType = 36
	FmpScriptSaveACopyAs                       FmpScriptStepType = 37
	FmpScriptOpenManageDatabase                FmpScriptStepType = 38
	FmpScriptSortRecords                       FmpScriptStepType = 39
	FmpScriptRelookupFieldContents             FmpScriptStepType = 40
	FmpScriptEnterPreviewMode                  FmpScriptStepType = 41
	FmpScriptPrintSetup                        FmpScriptStepType = 42
	FmpScriptPrint                             FmpScriptStepType = 43
	FmpScriptExitApplication                   FmpScriptStepType = 44
	FmpScriptUndoRedo                          FmpScriptStepType = 45
	FmpScriptCut                               FmpScriptStepType = 46
	FmpScriptCopy                              FmpScriptStepType = 47
	FmpScriptPaste                             FmpScriptStepType = 48
	FmpScriptClear                             FmpScriptStepType = 49
	FmpScriptSelectAll                         FmpScriptStepType = 50
	FmpScriptRevertRecordRequest               FmpScriptStepType = 51
	FmpScriptEnterBrowserMode                  FmpScriptStepType = 55
	FmpScriptInsertPicture                     FmpScriptStepType = 56
	FmpScriptSendEvent                         FmpScriptStepType = 57
	FmpScriptInsertCurrentUserName             FmpScriptStepType = 60
	FmpScriptInsertText                        FmpScriptStepType = 61
	FmpScriptPauseResumeScript                 FmpScriptStepType = 62
	FmpScriptSendMail                          FmpScriptStepType = 63
	FmpScriptSendDDEExecute                    FmpScriptStepType = 64
	FmpScriptDialPhone                         FmpScriptStepType = 65
	FmpScriptSpeak                             FmpScriptStepType = 66
	FmpScriptPerformApplescript                FmpScriptStepType = 67
	FmpScriptIf                                FmpScriptStepType = 68
	FmpScriptElse                              FmpScriptStepType = 69
	FmpScriptEndIf                             FmpScriptStepType = 70
	FmpScriptLoop                              FmpScriptStepType = 71
	FmpScriptExitLoopIf                        FmpScriptStepType = 72
	FmpScriptEndLoop                           FmpScriptStepType = 73
	FmpScriptGoToRelatedRecord                 FmpScriptStepType = 74
	FmpScriptCommitRecordsRequests             FmpScriptStepType = 75
	FmpScriptSetField                          FmpScriptStepType = 76
	FmpScriptInsertCalculatedResult            FmpScriptStepType = 77
	FmpScriptFreezeWindow                      FmpScriptStepType = 79
	FmpScriptRefreshWindow                     FmpScriptStepType = 80
	FmpScriptScrollWindow                      FmpScriptStepType = 81
	FmpScriptNewFile                           FmpScriptStepType = 82
	FmpScriptChangePassword                    FmpScriptStepType = 83
	FmpScriptSetMultiUser                      FmpScriptStepType = 84
	FmpScriptAllowUserAbort                    FmpScriptStepType = 85
	FmpScriptSetErrorCapture                   FmpScriptStepType = 86
	FmpScriptShowCustomDialog                  FmpScriptStepType = 87
	FmpScriptOpenScriptWorkspace               FmpScriptStepType = 88
	FmpScriptBlankLineComment                  FmpScriptStepType = 89
	FmpScriptHaltScript                        FmpScriptStepType = 90
	FmpScriptReplaceFieldContents              FmpScriptStepType = 91
	FmpScriptShowHideTextRuler                 FmpScriptStepType = 92
	FmpScriptBeep                              FmpScriptStepType = 93
	FmpScriptSetUseSystemFormats               FmpScriptStepType = 94
	FmpScriptRecoverFile                       FmpScriptStepType = 95
	FmpScriptSaveACopyAsAddOnPackage           FmpScriptStepType = 96
	FmpScriptSetZoomLevel                      FmpScriptStepType = 97
	FmpScriptCopyAllRecordsRequests            FmpScriptStepType = 98
	FmpScriptGoToPortalRow                     FmpScriptStepType = 99
	FmpScriptCopyRecordRequest                 FmpScriptStepType = 101
	FmpScriptFluchCacheToDisk                  FmpScriptStepType = 102
	FmpScriptExitScript                        FmpScriptStepType = 103
	FmpScriptDeletePortalRow                   FmpScriptStepType = 104
	FmpScriptOpenPreferences                   FmpScriptStepType = 105
	FmpScriptCorrectWord                       FmpScriptStepType = 106
	FmpScriptSpellingOptions                   FmpScriptStepType = 107
	FmpScriptSelectDictionaries                FmpScriptStepType = 108
	FmpScriptEditUserDictionary                FmpScriptStepType = 109
	FmpScriptOpenUrl                           FmpScriptStepType = 111
	FmpScriptOpenManageValueLists              FmpScriptStepType = 112
	FmpScriptOpenSharing                       FmpScriptStepType = 113
	FmpScriptOpenFileOptions                   FmpScriptStepType = 114
	FmpScriptAllowFormattingBar                FmpScriptStepType = 115
	FmpScriptSetNextSerialValue                FmpScriptStepType = 116
	FmpScriptExecuteSQL                        FmpScriptStepType = 117
	FmpScriptOpenHosts                         FmpScriptStepType = 118
	FmpScriptMoveResizeWindow                  FmpScriptStepType = 119
	FmpScriptArrangeAllWindows                 FmpScriptStepType = 120
	FmpScriptCloseWindow                       FmpScriptStepType = 121
	FmpScriptNewWindow                         FmpScriptStepType = 122
	FmpScriptSelectWindow                      FmpScriptStepType = 123
	FmpScriptSetWindowTitle                    FmpScriptStepType = 124
	FmpScriptElseIf                            FmpScriptStepType = 125
	FmpScriptConstrainFoundSet                 FmpScriptStepType = 126
	FmpScriptExtendFoundSet                    FmpScriptStepType = 127
	FmpScriptPerformFindReplace                FmpScriptStepType = 128
	FmpScriptOpenFindReplace                   FmpScriptStepType = 129
	FmpScriptSetSelection                      FmpScriptStepType = 130
	FmpScriptInsertFile                        FmpScriptStepType = 131
	FmpScriptExportFieldContents               FmpScriptStepType = 132
	FmpScriptOpenRecordRequest                 FmpScriptStepType = 133
	FmpScriptAddAccount                        FmpScriptStepType = 134
	FmpScriptDeleteAccount                     FmpScriptStepType = 135
	FmpScriptResetAccountPassword              FmpScriptStepType = 136
	FmpScriptEnableAccount                     FmpScriptStepType = 137
	FmpScriptRelogin                           FmpScriptStepType = 138
	FmpScriptConvertFile                       FmpScriptStepType = 139
	FmpScriptOpenManageDataSources             FmpScriptStepType = 140
	FmpScriptSetVariable                       FmpScriptStepType = 141
	FmpScriptInstallMenuSet                    FmpScriptStepType = 142
	FmpScriptSaveRecordsAsExcel                FmpScriptStepType = 143
	FmpScriptSaveRecordsAsPdf                  FmpScriptStepType = 144
	FmpScriptGoToObject                        FmpScriptStepType = 145
	FmpScriptSetWebViewer                      FmpScriptStepType = 146
	FmpScriptSetFieldByName                    FmpScriptStepType = 147
	FmpScriptInstallOntimerScript              FmpScriptStepType = 148
	FmpScriptOpenEditSavedFinds                FmpScriptStepType = 149
	FmpScriptPerformQuickFind                  FmpScriptStepType = 150
	FmpScriptOpenManageLayouts                 FmpScriptStepType = 151
	FmpScriptSaveRecordsAsSnapshotLink         FmpScriptStepType = 152
	FmpScriptSortRecordsByField                FmpScriptStepType = 154
	FmpScriptFindMatchingRecords               FmpScriptStepType = 155
	FmpScriptManageContainers                  FmpScriptStepType = 156
	FmpScriptInstallPluginFile                 FmpScriptStepType = 157
	FmpScriptInsertPdf                         FmpScriptStepType = 158
	FmpScriptInsertAudioVideo                  FmpScriptStepType = 159
	FmpScriptInsertFromUrl                     FmpScriptStepType = 160
	FmpScriptInsertFromDevice                  FmpScriptStepType = 161
	FmpScriptPerformScriptOnServer             FmpScriptStepType = 164
	FmpScriptOpenManageThemes                  FmpScriptStepType = 165
	FmpScriptShowHideMenubar                   FmpScriptStepType = 166
	FmpScriptRefreshObject                     FmpScriptStepType = 167
	FmpScriptSetLayoutObjectAnimation          FmpScriptStepType = 168
	FmpScriptClosePopover                      FmpScriptStepType = 169
	FmpScriptOpenUploadToHost                  FmpScriptStepType = 172
	FmpScriptEnableTouchKeyboard               FmpScriptStepType = 174
	FmpScriptPerformJavascriptInWebViewer      FmpScriptStepType = 175
	FmpScriptCommentedOut                      FmpScriptStepType = 176
	FmpScriptAvplayerPlay                      FmpScriptStepType = 177
	FmpScriptAvplayerSetPlaybackState          FmpScriptStepType = 178
	FmpScriptAvplayerSetOptions                FmpScriptStepType = 179
	FmpScriptRefreshPortal                     FmpScriptStepType = 180
	FmpScriptGetFolderPath                     FmpScriptStepType = 181
	FmpScriptTruncateTable                     FmpScriptStepType = 182
	FmpScriptOpenFavorites                     FmpScriptStepType = 183
	FmpScriptConfigureRegionMonitorScript      FmpScriptStepType = 185
	FmpScriptConfigureLocalNotification        FmpScriptStepType = 187
	FmpScriptGetFileExists                     FmpScriptStepType = 188
	FmpScriptGetFileSize                       FmpScriptStepType = 189
	FmpScriptCreateDataFile                    FmpScriptStepType = 190
	FmpScriptOpenDataFile                      FmpScriptStepType = 191
	FmpScriptWriteToDataFile                   FmpScriptStepType = 192
	FmpScriptReadFromDataFile                  FmpScriptStepType = 193
	FmpScriptGetDataFilePosition               FmpScriptStepType = 194
	FmpScriptSetDataFilePosition               FmpScriptStepType = 195
	FmpScriptCloseDataFile                     FmpScriptStepType = 196
	FmpScriptDeleteFile                        FmpScriptStepType = 197
	FmpScriptRenameFile                        FmpScriptStepType = 199
	FmpScriptSetErrorLogging                   FmpScriptStepType = 200
	FmpScriptConfigureNfcReading               FmpScriptStepType = 201
	FmpScriptConfigureMachineLearningModel     FmpScriptStepType = 202
	FmpScriptExecuteFileMakerDataAPI           FmpScriptStepType = 203
	FmpScriptOpenTransaction                   FmpScriptStepType = 205
	FmpScriptCommitTransaction                 FmpScriptStepType = 206
	FmpScriptRevertTransaction                 FmpScriptStepType = 207
	FmpScriptSetSessionIdentifier              FmpScriptStepType = 208
	FmpScriptSetDictionary                     FmpScriptStepType = 209
	FmpScriptPerformScriptOnServerWithCallback FmpScriptStepType = 210
	FmpScriptTriggerClarisConnectFlow          FmpScriptStepType = 211
	FmpScriptAssert                            FmpScriptStepType = 255
)
