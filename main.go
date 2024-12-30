package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// Listen on UDP port 13213 on all interfaces
	addr := "0.0.0.0:13213"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		fmt.Printf("Failed to start UDP server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("[%s] UDP server listening on %s\n", time.Now().Format(time.RFC3339), addr)

	buffer := make([]byte, 1024) // Buffer to store incoming data
	for {
		n, remoteAddr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Printf("[%s] Error reading from connection: %v\n", time.Now().Format(time.RFC3339), err)
			continue
		}

		// Get the current timestamp
		timestamp := time.Now().Format(time.RFC3339)

		// Print received hex data
		hexString := "0x"
		for _, b := range buffer[:n] {
			hexString += fmt.Sprintf("%02X", b)
		}
		fmt.Printf("[%s] Received message from %s: %s\n", timestamp, remoteAddr, hexString)

		// Ensure we have enough data (at least 25 bytes to extract the 25th byte)
		if n >= 25 {
			// Parse AVL Data Packet
			// First 6 bytes and 25th byte for the response
			response := append(buffer[0:6], buffer[24:25]...)

			// Parse AVL Data
			imei := string(buffer[6:23]) // IMEI from byte 6 to 21
			codecID := fmt.Sprintf("%02X", buffer[23]) // Codec ID at byte 21
			timestampHex := ""
			for _, b := range buffer[24:32] {
				timestampHex += fmt.Sprintf("%02X", b)
			}
			numberdata := fmt.Sprintf("%02X", buffer[24]) // Priority at byte 31

			// Logging the parsed AVL Data
			fmt.Printf("[%s] Parsed AVL Data Packet:\n", timestamp)
			fmt.Printf("  IMEI: %s\n", imei)
			fmt.Printf("  Codec ID: %s\n", codecID)
			fmt.Printf("  Timestamp: %s\n", timestampHex)
			fmt.Printf("  numberof data: %s\n", numberdata)

			// Construct the full response: First 6 bytes + 25th byte
			// Convert the response to hexadecimal format for logging
			responseHex := "0x"
			for _, b := range response {
				responseHex += fmt.Sprintf("%02X", b)
			}

			// Send the response back to the client
			_, err := conn.WriteTo(response, remoteAddr)
			if err != nil {
				fmt.Printf("[%s] Error sending reply: %v\n", timestamp, err)
			} else {
				fmt.Printf("[%s] Sent reply to %s: %s\n", timestamp, remoteAddr, responseHex)
			}
		} else {
			// If the message is too short to parse
			fmt.Printf("[%s] Received packet is too short to form a response.\n", timestamp)
		}
	}
}
