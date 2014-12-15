/*
The MIT License (MIT)

Copyright (c) 2014 DutchCoders [https://github.com/dutchcoders/]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import ()

type RuleState int

const (
	RuleStateWhitelist       RuleState = 1
	RuleStateBlacklist       RuleState = 2
	RuleStateSilentBlacklist RuleState = 3
	RuleStateRemove          RuleState = 4
)

type RuleType int

const (
	RuleTypeBinary      RuleType = 1
	RuleTypeCertificate RuleType = 2
)

type ClientMode int

const (
	ClientModeMonitor  ClientMode = 1
	ClientModeLockdown ClientMode = 2
)

type EventState int

const (
	EventStateAllowMonitor     EventState = 1
	EventStateAllowBinary      EventState = 2
	EventStateAllowCertificate EventState = 3
	EventStateBlockUnknown     EventState = 4
	EventStateBlockBinary      EventState = 5
	EventStateBlockCertificate EventState = 6
)
