package jsonselect

import (
    "fmt"
    "log"
    "os"
    "strings"
    "github.com/latestrevision/go-simplejson"
)

type Logger struct {
    Enabled bool
    recursionLevel int
    prefixes map[int]string
}


var logger = Logger{false, 0, nil}
var handler = log.New(os.Stderr, "jsonselect: ", 0)
var recursionMarker = "⇢ "

func (l *Logger) formatPrefix(a ...interface{}) []interface{} {
    var arguments []interface{}
    arguments = append(
        arguments,
        strings.Repeat(recursionMarker, l.recursionLevel),
    )
    prefix, ok := l.prefixes[l.recursionLevel]
    if ok {
        arguments = append(
            arguments,
            prefix,
        )
    }
    arguments = append(
        arguments,
        a...,
    )
    return arguments
}


func (l *Logger) Print(a ...interface{}) {
    if logger.Enabled {
        handler.Print(l.formatPrefix(a...)...)
    }
}

func (l *Logger) Println(a ...interface{}) {
    if logger.Enabled {
        handler.Println(l.formatPrefix(a...)...)
    }
}

func (l *Logger) IncreaseDepth() {
    if logger.Enabled {
        l.recursionLevel++
    }
}

func (l *Logger) DecreaseDepth() {
    if logger.Enabled {
        l.recursionLevel--
    }
}

func (l *Logger) SetPrefix(prefix ...interface{}) {
    if logger.Enabled {
        l.prefixes[l.recursionLevel] = fmt.Sprint(prefix...)
    }
}

func (l *Logger) ClearPrefix() {
    if logger.Enabled {
        l.prefixes[l.recursionLevel] = ""
    }
}

func EnableLogger() {
    logger.prefixes = make(map[int]string)
    logger.Enabled = true
}

func getFormattedNodeMap(nodes map[*simplejson.Json]*jsonNode) []string {
    output := make([]*jsonNode, 0, len(nodes))
    for _, val := range nodes {
        output = append(output, val)
    }
    return getFormattedNodeArray(output)
}

func getFormattedNodeArray(nodes []*jsonNode) []string {
    var formatted []string
    for _, node := range nodes {
        if node != nil {
            formatted = append(formatted, fmt.Sprint(*node))
        } else {
            formatted = append(formatted, fmt.Sprint(nil))
        }
    }
    return formatted
}

func getFormattedTokens(tokens []*token) []string {
    var output []string
    for _, token := range tokens {
        output = append(output, fmt.Sprint(token.val))
    }
    return output
}

func getFormattedExpression(tokens []*exprElement) []string {
    var output []string
    for _, token := range tokens {
        output = append(output, fmt.Sprint(token.value))
    }
    return output
}
