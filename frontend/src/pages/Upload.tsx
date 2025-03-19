import UploadForm from "../components/UploadForm";

export default function Upload() {
  return (
    <div className="flex flex-col items-center">
      <h1 className="text-4xl font-bold my-4">Welcome to My Portfolio</h1>
      <UploadForm />
    </div>
  )
}