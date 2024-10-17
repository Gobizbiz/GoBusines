<?php
// Mulai sesi
session_start();

// Ambil data dari form login
$username = $_POST['username'];
$password = $_POST['password'];

// Koneksi ke database MySQL
$conn = new mysqli('localhost', 'root', '', 'login_db');

// Periksa koneksi
if ($conn->connect_error) {
    die("Koneksi gagal: " . $conn->connect_error);
}

// Query untuk memeriksa username
$stmt = $conn->prepare("SELECT * FROM users WHERE username = ?");
$stmt->bind_param("s", $username);
$stmt->execute();
$result = $stmt->get_result();

if ($result->num_rows > 0) {
    $row = $result->fetch_assoc();
    // Verifikasi password
    if (password_verify($password, $row['password'])) {
        $_SESSION['username'] = $username;
        header("Location: dashboard.php"); // Redirect ke dashboard jika login berhasil
    } else {
        echo "Password salah!";
    }
} else {
    echo "Username tidak ditemukan!";
}

$stmt->close();
$conn->close();
?>
