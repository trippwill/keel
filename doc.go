// Package keel is a layout manager for terminal applications.
//
// Box Model:
//
//	+----------------+
//	|    Chrome      |
//	| +------------+ |
//	| |  Content   | |
//	| +------------+ |
//	+----------------+
//
// Chrome: Borders, padding, margins
// Content: The area where child elements are rendered
// Cell: A single character position in the terminal
package keel
