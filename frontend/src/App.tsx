import ImageList from "./components/ImageList";
import UploadForm from "./components/UploadForm";

function App() {
  return (
    <div className="min-h-screen items-center justify-center bg-gray-200">
      <UploadForm />
      <ImageList />
    </div>
  )
}

export default App;