package store

// We can always have similar methods in inMemory without using storage (rsync does it this way)
// we can on real time/per request find the checksum of each file to see if a same file exist witha  different name
// But that will not be scalable for even a small amount of files.

// Add
// GetFileByName
// GetFileBySameData
// Update
// Delete
