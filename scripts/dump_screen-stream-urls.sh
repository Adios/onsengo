#!/bin/bash
../onsengo dump | sed -E 's/"streaming_url":"[^"]+"/"streaming_url":"HAS_BEEN_SCREENED"/g'
