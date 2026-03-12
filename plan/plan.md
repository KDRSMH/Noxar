## Proje Planı: GoRecon CLI Aracı Tasarımı ve Uygulaması

TL;DR – Goroutineler ve kanallar gibi Go’nun eşzamanlılık özelliklerini gösteren, `GoRecon` adlı ağ keşif (network reconnaissance) CLI aracını tanımlayıp inşa et. Proje farklı keşif tekniklerini içeren modüler paketler, bir CLI arayüzü ve kod içinde ayrıntılı dokümantasyon ile test rehberi barındıracak. Plan, yapı, algoritmalar, bağımlılıklar ve öğrenme hedeflerini kapsar; uygulama bu plana göre ilerleyecek.

**Adımlar**

1. **Kapsam ve özellik tanımı** (Keşif/analiz): araç ne yapıyor, gerçek dünyada nerede kullanılır, hangi keşif tekniklerini destekliyor – “Projenin Tanımı” bölümü için kullanılacak.
2. **Dizin/yapı tasarımı**: her klasör/dosya için kısa açıklama içeren bir hiyerarşi hazırlamak. `cmd/`, `pkg/hosts`, `pkg/ports`, `pkg/dns`, `pkg/utils` gibi paketleri belirlemek.
3. **Modül mimarisi ve algoritma akışı**: her paket için sorumluluk ve adım adım davranışını tanımlamak. Goroutineler/kanallar nerede kullanılır belirlemek. Sonuçlar için veri yapıları ve modüller arası iletişimi tasarlamak.
4. **Bağımlılıkları seçmek**: standart kütüphane paketleri (`net`, `fmt`, `os`, `sync`, `flag`/`cobra`) ve varsa üçüncü taraf paketler (`github.com/spf13/cobra` CLI için, `github.com/miekg/dns` DNS için) seçmek. `go.mod` başlangıç içeriğini hazırlamak.
5. **Kod yazmak**: modülleri sıralı veya paralel gruplar halinde uygulamak. Her dosyada ayrıntılı yorumlar, fonksiyon dokümantasyonu ve eşzamanlılık açıklamaları yer alacak. Her büyük özellik için bir dosya hedeflenmeli.
6. **CLI’yi uygulamak**: `main.go` (veya `cmd/` altındaki komut) bayrakları ayrıştıracak, modül çağrılarını düzenleyecek, goroutineleri yönetecek ve kanallar aracılığıyla sonuçları toplayacak.
7. **Testler eklemek**: parser’lar ve tarama işlevleri için birim testleri ve localhost/denetimli ortam karşısında birkaç entegrasyon senaryosu yazmak. Güvenli test hedeflerini belgelendirmek.
8. **README/Dokümantasyon yazmak**: CLI örnekleri, bayrak açıklamaları, çıktı örnekleri ve kullanım adım‑adım rehberi eklemek. Her modülde öğrenilen Go kavramları ve haftalık öğrenme hedefleri bölümü eklemek.
9. **Doğrulama**: `go build`, `go test ./...` çalıştırmak, elle kullanım örnekleri denemek. Eşzamanlılık için `-race` dedektörünü kullanmak. Doğrulama adımlarını proje dokümanına eklemek.

**İlgili dosyalar**

- `go.mod` — modül bildirimi ve bağımlılıklar  
- `cmd/gorecon/main.go` — cobra veya flag kullanan giriş noktası ve CLI mantığı  
- `pkg/hosts/hosts.go` — ana bilgisayar (host) keşif fonksiyonları  
- `pkg/ports/scan.go` — port tarama yardımcıları  
- `pkg/dns/dns.go` — ters/doğrudan DNS sorguları  
- `pkg/utils/results.go` — paylaşılan sonuç yapıları ve kanal tanımları  
- `README.md` — proje genel bakışı ve kullanım  
- `pkg/.../*_test.go` — her paket için birim testler

**Doğrulama**

1. `go build ./cmd/gorecon` ile derleyip binary’nin çalıştığından emin olmak.  
2. `127.0.0.1` ve güvenli bir alan adına karşı örnek keşif komutları çalıştırıp çıktıları beklenenle karşılaştırmak.  
3. `go test -race ./...` ile testleri çalıştırıp eşzamanlılık sorunlarını yakalamak.  
4. CLI yardım ve bayrak ayrıştırmanın belgelenildiği gibi çalıştığını doğrulamak.

**Kararlar**

- Paralel ana bilgisayar/port kontrolleri için goroutineler; sonuçları toplamak ve tamamlanmayı bildirmek için kanallar kullanılacak.  
- CLI bayrakları ve alt komutlar için cobra tercih edilebilir (Go projelerinde yaygın), flag daha basit ama öğrenme hedefine göre seçilecek.  
- Kapsam yalnızca ağ keşfiyle sınırlı; sömürü veya saldırı içermeyecek.  
- Harici veritabanı yok; sonuçlar stdout’a veya isteğe bağlı JSON dosyasına yazılacak.

**Ek Hususlar**

1. Erken aşamada kabul edilebilir üçüncü taraf paketleri belirle; sadelik istenirse sadece standart kütüphane kullan ve DNS/CLI işlevlerini kendin yaz.  
2. DNS sorguları ile port taramayı eşzamanlı çalıştırıp merkezî bir koordine ile mi yoksa kanallarla zincirleyerek mi iletişim sağlayacağını netleştir.  
3. Çıktı formatını (düz metin vs. JSON) ve bir ayrıntı/verbose bayrağı eklenip eklenmeyeceğini kararlaştır.

---

Elindeki bu plan, `GoRecon` projesinin inşası için yol haritasını ve istenen dokümantasyonu içerir. İnceledikten sonra, belirli bölümleri detaylandırmamı veya kod ve Markdown çıktısını oluşturmamı isteyebilirsin.
