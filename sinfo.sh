#!/bin/bash

BINARY_PATH="~/scripts/shour/shour"
PROGRESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoibnNoZXJpIiwiZXhwIjoxNzI5MjQzNDgwfQ.SLSLEt_KuLppJOvNVOT7hoP1nK7uUg1mtjpSBbNGFaY"
PLATFORM_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwIjp7InVpZCI6IjgxOTQxZDVhLWEyNTUtNDRlNi04NjBlLTk3MmI2ZmM1OTQwYiIsInIiOiJzdHVkZW50In0sImV4cCI6MTcyOTI0MzQ3Nn0.gS2kvIPSQakEnopBf26sLAh9PFJCaKzR0tCBhDwSU-A"
HOURS_REQUIREMENT=20 # 20 OR 30

hoursText="Hours fulfilled!"
reviewsText="You have upcoming reviews."


sound_bell="/usr/share/sounds/freedesktop/stereo/bell.oga"

isReviewsWorked=false


gnome-terminal --geometry=36x11 -- bash -c "
    while true; do
        clear
        
        $BINARY_PATH $PROGRESS_TOKEN $PLATFORM_TOKEN $HOURS_REQUIREMENT

        sleep 60  # 30 secs
    done
"

: << 'END_COMMENT'
        if $BINARY_PATH $PROGRESS_TOKEN $PLATFORM_TOKEN | grep -q \"No upcoming reviews\" && [ \"\$isReviewsWorked\" = false ]; then
            paplay \"$sound_bell\"
            notify-send \"Student Info (Sinfo)\" \"You have upcoming reviews.\"
            isReviewsWorked=true
        fi
END_COMMENT
