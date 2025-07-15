package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Barang struct {
	ID         int
	KodeBarang string
	NamaBarang string
	StokBarang int
	HargaBeli  float64
	HargaJual  float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "./AplikasiGudang.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Buat table jika belum ada
	sqlTable := `CREATE TABLE IF NOT EXISTS AplikasiGudang (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		kodeBarang TEXT,
		namaBarang TEXT,
		stokBarang INTEGER,
		hargaBeli REAL,
		hargaJual REAL,
		CreatedAt TIMESTAMP,
		UpdatedAt TIMESTAMP
	);`

	_, err = db.Exec(sqlTable)
	if err != nil {
		panic(err)
	}

	var pilihan int
	var ulang string

	for {
		fmt.Println("\n=== Aplikasi Manajemen Gudang ===")
		fmt.Println("1. Tambah Barang")
		fmt.Println("2. List Barang")
		fmt.Println("3. Edit Barang")
		fmt.Println("4. Hapus Barang")
		fmt.Println("5. Keluar")
		fmt.Print("Pilih menu: ")
		fmt.Scanln(&pilihan)

		switch pilihan {
		case 1:
			tambahBarang()
		case 2:
			listBarang()
		case 3:
			editBarang()
		case 4:
			hapusBarang()
		case 5:
			fmt.Println("Terima kasih telah menggunakan aplikasi ini.")
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid!")
			fmt.Print("Ingin memulai ulang? (y/n): ")
			fmt.Scanln(&ulang)
			if ulang == "n" || ulang == "N" {
				fmt.Println("Terima kasih telah menggunakan aplikasi ini.")
				os.Exit(0)
			}
		}
	}
}

func tambahBarang() {
	var jumlah int
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan jumlah barang yang ingin ditambahkan: ")
	for {
		_, err := fmt.Scanln(&jumlah)
		if err != nil || jumlah <= 0 {
			fmt.Println("Input harus berupa angka positif, coba lagi.")
			var dummy string
			fmt.Scanln(&dummy)
			fmt.Print("Masukkan jumlah barang yang ingin ditambahkan: ")
			continue
		}
		break
	}

	for i := 0; i < jumlah; i++ {
		fmt.Printf("\n=== Barang ke-%d ===\n", i+1)

		fmt.Print("Kode Barang: ")
		var kodeBarang string
		fmt.Scanln(&kodeBarang)

		fmt.Print("Nama Barang: ")
		namaBarang, _ := reader.ReadString('\n')
		namaBarang = strings.TrimSpace(namaBarang)

		var stokBarang int
		fmt.Print("Stok Barang: ")
		for {
			_, err := fmt.Scanln(&stokBarang)
			if err != nil {
				fmt.Println("Input harus angka, coba lagi.")
				var dummy string
				fmt.Scanln(&dummy)
				fmt.Print("Stok Barang: ")
				continue
			}
			break
		}

		var hargaBeli float64
		fmt.Print("Harga Beli: ")
		for {
			_, err := fmt.Scanln(&hargaBeli)
			if err != nil {
				fmt.Println("Input harus angka, coba lagi.")
				var dummy string
				fmt.Scanln(&dummy)
				fmt.Print("Harga Beli: ")
				continue
			}
			break
		}

		var hargaJual float64
		fmt.Print("Harga Jual: ")
		for {
			_, err := fmt.Scanln(&hargaJual)
			if err != nil {
				fmt.Println("Input harus angka, coba lagi.")
				var dummy string
				fmt.Scanln(&dummy)
				fmt.Print("Harga Jual: ")
				continue
			}
			break
		}

		_, err := db.Exec(`
			INSERT INTO AplikasiGudang
			(kodeBarang, namaBarang, stokBarang, hargaBeli, hargaJual, CreatedAt, UpdatedAt)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			kodeBarang, namaBarang, stokBarang, hargaBeli, hargaJual, time.Now(), time.Now(),
		)
		if err != nil {
			fmt.Println("Gagal menyimpan ke database:", err)
		} else {
			fmt.Println("Barang berhasil disimpan ke database.")
		}
	}
}

func listBarang() {
	rows, err := db.Query(`SELECT id, kodeBarang, namaBarang, stokBarang, hargaBeli, hargaJual, CreatedAt FROM AplikasiGudang`)
	if err != nil {
		fmt.Println("Gagal mengambil data:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n=== Daftar Barang ===")
	for rows.Next() {
		var b Barang
		err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.StokBarang, &b.HargaBeli, &b.HargaJual, &b.CreatedAt)
		if err != nil {
			fmt.Println("Error membaca data:", err)
			continue
		}
		fmt.Printf("ID: %d | Kode: %s | Nama: %s | Stok: %d | Harga Beli: %.2f | Harga Jual: %.2f | CreatedAt: %v\n",
			b.ID, b.KodeBarang, b.NamaBarang, b.StokBarang, b.HargaBeli, b.HargaJual, b.CreatedAt)
	}
}

func hapusBarang() {
	var id int
	fmt.Print("Masukkan ID Barang yang ingin dihapus: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Input tidak valid.")
		return
	}

	result, err := db.Exec("DELETE FROM AplikasiGudang WHERE id = ?", id)
	if err != nil {
		fmt.Println("Gagal menghapus barang:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("Data dengan ID tersebut tidak ditemukan.")
	} else {
		fmt.Println("Barang berhasil dihapus.")
	}
}

func editBarang() {
	var id int
	fmt.Print("Masukkan ID Barang yang ingin diedit: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Input tidak valid.")
		return
	}

	// Ambil data lama
	var kodeLama, namaLama string
	var stokLama int
	var hargaBeliLama, hargaJualLama float64
	var createdAt, updatedAt time.Time

	row := db.QueryRow(`SELECT kodeBarang, namaBarang, stokBarang, hargaBeli, hargaJual, CreatedAt, UpdatedAt 
		FROM AplikasiGudang WHERE id = ?`, id)

	err = row.Scan(&kodeLama, &namaLama, &stokLama, &hargaBeliLama, &hargaJualLama, &createdAt, &updatedAt)
	if err != nil {
		fmt.Println("Data tidak ditemukan:", err)
		return
	}

	fmt.Println("\n=== Data Lama ===")
	fmt.Printf("Kode: %s\nNama: %s\nStok: %d\nHarga Beli: %.2f\nHarga Jual: %.2f\n",
		kodeLama, namaLama, stokLama, hargaBeliLama, hargaJualLama)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan Nama Baru (Enter untuk tidak diubah): ")
	namaBaru, _ := reader.ReadString('\n')
	namaBaru = strings.TrimSpace(namaBaru)
	if namaBaru == "" {
		namaBaru = namaLama
	}

	fmt.Print("Masukkan Stok Baru (Enter untuk tidak diubah): ")
	var inputStok string
	stokBaru := stokLama
	fmt.Scanln(&inputStok)
	if inputStok != "" {
		fmt.Sscanf(inputStok, "%d", &stokBaru)
	}

	fmt.Print("Masukkan Harga Beli Baru (Enter untuk tidak diubah): ")
	var inputHargaBeli string
	hargaBeliBaru := hargaBeliLama
	fmt.Scanln(&inputHargaBeli)
	if inputHargaBeli != "" {
		fmt.Sscanf(inputHargaBeli, "%f", &hargaBeliBaru)
	}

	fmt.Print("Masukkan Harga Jual Baru (Enter untuk tidak diubah): ")
	var inputHargaJual string
	hargaJualBaru := hargaJualLama
	fmt.Scanln(&inputHargaJual)
	if inputHargaJual != "" {
		fmt.Sscanf(inputHargaJual, "%f", &hargaJualBaru)
	}

	_, err = db.Exec(`UPDATE AplikasiGudang
		SET namaBarang = ?, stokBarang = ?, hargaBeli = ?, hargaJual = ?, UpdatedAt = ?
		WHERE id = ?`,
		namaBaru, stokBaru, hargaBeliBaru, hargaJualBaru, time.Now(), id)
	if err != nil {
		fmt.Println("Gagal mengupdate data:", err)
		return
	}

	fmt.Println("Data berhasil diperbarui.")
}
