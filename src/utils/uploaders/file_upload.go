package uploaders

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Fungsi untuk meng-upload file ke folder storage/files dan mengembalikan path relatif
// func UploadFile(ctx context.Context, folder string) (string, error) {
// 	// Ambil request dari context
// 	req := ctx.Value("request").(*http.Request)
// 	if req == nil {
// 		return "", errors.New("request not found in context")
// 	}

// 	// Ambil file dari form-data
// 	file, fileHeader, err := req.FormFile("file")
// 	if err != nil {
// 		return "", errors.Wrap(err, "failed to get file from form-data")
// 	}

// 	// Buat folder storage/files jika belum ada
// 	storagePath := "storage/files"
// 	folderPath := filepath.Join(storagePath, folder)
// 	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
// 		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
// 			return "", errors.Wrap(err, "failed to create folder")
// 		}
// 	}

// 	// Generate nama file unik
// 	fileExtension := filepath.Ext(fileHeader.Filename)
// 	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)

// 	// Tentukan path penyimpanan
// 	filePath := filepath.Join(folderPath, fileName)

// 	// Simpan file ke folder yang ditentukan
// 	dst, err := os.Create(filePath)
// 	if err != nil {
// 		return "", errors.Wrap(err, "failed to create file")
// 	}
// 	defer dst.Close()

// 	// Salin konten file ke destination file
// 	if _, err := dst.ReadFrom(file); err != nil {
// 		return "", errors.Wrap(err, "failed to save file")
// 	}

// 	// Mengembalikan path relatif yang bisa disimpan di database
// 	return fmt.Sprintf("/files/%s/%s", folder, fileName), nil
// }

// Fungsi untuk meng-upload file ke folder storage/files dan mengembalikan path relatif
func UploadFile(req *http.Request, folder string) (string, error) {
	// Ambil file dari form-data
	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		return "", errors.Wrap(err, "failed to get file from form-data")
	}
	defer file.Close()

	// Buat folder storage/files jika belum ada
	storagePath := "src/storage/files"
	folderPath := filepath.Join(storagePath, folder)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return "", errors.Wrap(err, "failed to create folder")
		}
	}

	// Generate nama file unik
	fileExtension := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)

	// Tentukan path penyimpanan
	filePath := filepath.Join(folderPath, fileName)

	// Simpan file ke folder yang ditentukan
	dst, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrap(err, "failed to create file")
	}
	defer dst.Close()

	// Salin konten file ke destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", errors.Wrap(err, "failed to save file")
	}

	// Mengembalikan path relatif yang bisa disimpan di database
	return fmt.Sprintf("/files/%s/%s", folder, fileName), nil
}
