#!/bin/bash

BINARY_PATH="~/scripts/shour/shour"
PROGRESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoibnNoZXJpIiwiZXhwIjoxNzI4ODMxMzY4fQ.X0KmPJGaDfQ7mFoB_3I7k4RppxNcWTIcHFEwhThcyI8"
PLATFORM_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwIjp7InVpZCI6IjgxOTQxZDVhLWEyNTUtNDRlNi04NjBlLTk3MmI2ZmM1OTQwYiIsInIiOiJzdHVkZW50In0sImV4cCI6MTcyODcxNjYzN30.lqkipshn2gstkUe4OuTxB77JqYgs0oop0_Mmg1599VQ"

hoursText="Hours fulfilled!"
reviewsText="You have upcoming reviews."


sound_bell="/usr/share/sounds/freedesktop/stereo/bell.oga"

isReviewsWorked=false


gnome-terminal --geometry=36x11 -- bash -c "
    while true; do
        clear
        
        $BINARY_PATH $PROGRESS_TOKEN $PLATFORM_TOKEN
        

        sleep 30  # 300 секунд = 5 минут
    done
" &

sleep 1

TERMINAL_ID=$(xdotool getactivewindow)
xprop -id $TERMINAL_ID -f _NET_WM_STATE 32a -set _NET_WM_STATE "_NET_WM_STATE_ABOVE"

: << 'END_COMMENT'
        if $BINARY_PATH $PROGRESS_TOKEN $PLATFORM_TOKEN | grep -q \"No upcoming reviews\" && [ \"\$isReviewsWorked\" = false ]; then
            paplay \"$sound_bell\"
            notify-send \"Student Info (Sinfo)\" \"You have upcoming reviews.\"
            isReviewsWorked=true  # Убедитесь, что здесь нет пробелов вокруг '='
        fi
END_COMMENT
