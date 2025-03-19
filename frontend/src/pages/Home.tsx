import ImageList from "../components/ImageList";

export default function Home() {
  return (
    <div className="flex flex-col items-center">
      <h1 className="text-4xl font-bold my-4">Welcome to My Portfolio</h1>
      <ImageList />
    </div>
  )
}