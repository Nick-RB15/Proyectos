package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Middleware para registrar solicitudes
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

// Estructura para el contacto
type Contacto struct {
	Email            string
	Telefono         string // Opcional
	NombreEmpresa    string
	NombreReclutador string // Opcional
	DatosExtras      string
}

func main() {
	// Configuración del servidor
	const port = ":8080"
	const staticDir = "./Style"
	const datosArchivo = "datos.json"

	// Servir archivos estáticos (CSS, imágenes, etc)
	http.Handle("/Style/", http.StripPrefix("/Style/", http.FileServer(http.Dir(staticDir))))

	// Ruta para la página principal (perfil)
	http.HandleFunc("/", logRequest(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "profile.html")
	}))

	// Ruta para la página de descripción
	http.HandleFunc("/description", logRequest(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "description.html")
	}))

	// Ruta para la página de contacto
	http.HandleFunc("/contact", logRequest(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "contact.html")
	}))

	// Manejador de acciones de contacto
	http.HandleFunc("/accion", logRequest(func(w http.ResponseWriter, r *http.Request) {
		// Solo permitir método POST
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Decodificar el cuerpo de la solicitud
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
			return
		}

		// Crear un contacto con datos de ejemplo
		contacto := Contacto{
			Email:            data["email"],
			Telefono:         data["telefono"], // Opcional: dejar vacío si no se proporciona
			NombreEmpresa:    data["nombre-empresa"],
			NombreReclutador: data["nombre-reclutador"], // Opcional: dejar vacío si no se proporciona
			DatosExtras:      data["datos-extras"],
		}

		// Convertir el contacto a un texto para su validación
		textoContacto := fmt.Sprintf("Email: %s, Teléfono: %s, Nombre de la Empresa: %s, Nombre del Reclutador: %s, Datos Extras: %s",
			contacto.Email, contacto.Telefono, contacto.NombreEmpresa, contacto.NombreReclutador, contacto.DatosExtras)

		// Validar el texto del contacto
		if textoContacto == "" {
			http.Error(w, "No se pudo procesar el contacto debido a la falta de información.", http.StatusBadRequest)
			return
		}

		// Guardar los datos en un archivo
		archivo, err := os.OpenFile(datosArchivo, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error al guardar los datos del contacto.", http.StatusInternalServerError)
			return
		} else {
			_, err = archivo.WriteString(fmt.Sprintf("%s\n", textoContacto))
			if err != nil {
				log.Println(err)
				http.Error(w, "Error al guardar los datos del contacto.", http.StatusInternalServerError)
			}
			archivo.Close()
		}

		// Enviar respuesta JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Tu correo ha sido enviado exitosamente."})
	}))

	// Ruta para la página de servicios
	http.HandleFunc("/services", logRequest(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "service.html")
	}))

	// Ruta para la demo de API
	http.HandleFunc("/api/demo/api", logRequest(func(w http.ResponseWriter, r *http.Request) {
		// Solo permitir método GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Obtener una broma aleatoria de Chuck Norris
		resp, err := http.Get("https://api.chucknorris.io/jokes/random")
		if err != nil {
			http.Error(w, "Error al obtener la demo de API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var datos map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&datos); err != nil {
			http.Error(w, "Error al procesar la demo de API", http.StatusBadRequest)
			return
		}

		// Filtrar solo el campo "value"
		filteredData := map[string]interface{}{
			"value": datos["value"],
		}

		// Enviar la respuesta filtrada como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredData)
	}))

	// Ruta para la demo de API de servicios
	http.HandleFunc("/api/demo/service", logRequest(func(w http.ResponseWriter, r *http.Request) {
		// Solo permitir método GET
		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// Preparar datos de demostración de servicios
		demoData := map[string]interface{}{
			"message": "Demostración de servicios exitosa.",
			"services": []map[string]string{
				{
					"name":        "APIs y Microservicios",
					"description": "Desarrollo de APIs RESTful y microservicios usando Go, Node.js y Python. Integración con bases de datos y servicios externos.",
				},
				{
					"name":        "Gestión de Bases de Datos",
					"description": "Diseño e implementación de bases de datos SQL y NoSQL. Optimización de consultas y modelado de datos.",
				},
				{
					"name":        "Logging",
					"description": "Implementación de sistemas de registro para aplicaciones web y servicios.",
				},
			},
		}

		// Enviar la demo de servicios como respuesta JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(demoData)
	}))

	// Iniciar el servidor
	log.Printf("Servidor iniciando en http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
