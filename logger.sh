#!/bin/bash

LOGFILE=zshlog_$(hostname)_$(date +%Y-%m-%d_%H%M%S_%N.log)
[ ! -d $LOGDIR ] && mkdir -p $LOGDIR
tmux  set-option default-terminal "screen" \; \
pipe-pane        "cat >> $LOGDIR/$LOGFILE" \; \
display-message  "💾Started logging to $LOGDIR/$LOGFILE"
