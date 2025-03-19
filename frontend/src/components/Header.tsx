import { Link } from "react-router-dom";

export default function Header() {
  return (
    <header className="bg-white shadow-md py-4 px-6 flex justify-between">
      <h1 className="text-xl font-bold">Photo Portfolio</h1>
      <nav>
        <Link to="/" className="mx-2 text-blue-500">Home</Link>
        <Link to="/gallery" className="mx-2 text-blue-500">Gallery</Link>
        <Link to="/upload" className="mx-2 text-blue-500">Upload</Link>
      </nav>
    </header>
  )
}