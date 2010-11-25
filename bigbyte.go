// Copyright 2010 bjarneh@ifi.uio.no. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bigbyte

// Boyer Moore Horspool; backwards matching string algorithm.
// Can be used to search for substrings in large []byte arrays
// NOTE: not as quick as Boyer Moore in worst case scenarios
// but without the added complexity in terms of code, and
// without the startup problems that Boyer Moore has. I.e. this
// will work comparable to bytes.Index on short strings, unlike
// Boyer Moore with it's good suffix shift calculation at startup.
func IndexBMH(haystack, needle []byte) (int) {

    const UINT8_MAX = 256

    var badshift [UINT8_MAX]int
    var i, offset, scan, last, maxoffset int

    if haystack == nil || needle == nil {
        return -1
    }

    if len(haystack) < len(needle) {
        return -1
    }

    if len(needle) == 0 {
        return 0
    }

    if len(needle) == 1 {
        for i = 0; i < len(haystack); i++ {
            if needle[0] == haystack[i] {
                return i
            }
        }
        return -1
    }

    // the algorithm starts here, special cases above.

    last = len(needle) - 1
    maxoffset = len(haystack) - len(needle)

    // default shift => entire needle length
    for i = 0; i < UINT8_MAX; i++ {
        badshift[i] = len(needle)
    }

    // override with length to next match in needle
    for i = 0; i < last; i++ {
        badshift[needle[i]] = last - i
    }

    for offset <= maxoffset {

        // start at end of pattern and match backwards
        for scan = last; needle[scan] == haystack[scan+offset]; scan-- {

            // we have a match
            if(scan == 0){
                return offset
            }

        }
        // use last character and badshift to calculate offset
        offset += badshift[haystack[offset + last]];
    }

    // indicates no match
    return -1
}
