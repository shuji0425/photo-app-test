import { Image } from "../types/image";

// サイレンダリング削減
const ImageCard = ({ image, priority }: { image: Image; priority?: boolean }) => {
  return (
    <div className="bg-white shadow-lg rounded-lg overflow-hidden transform hover:scale-105 transition-transform duration-300">
      <div className="relative w-full" style={{ minHeight: "300px"}}>
        <img
          srcSet={`
            ${image.url}?w=400&h=300&fit=crop&fm=webp 400w,
            ${image.url}?w=800&h=600&fit=crop&fm=webp 800w,
            ${image.url}?w=1200&h=900&fit=crop&fm=webp 1200w
          `}
          sizes="(max-width: 640px) 400px, (max-width: 1024px) 800px, 1200px"
          src={`${image.url}?w=800&h=600&fit=crop&fm=webp`}
          alt={`Image ${image.id}`}
          width={800}
          height={600}
          className="w-full h-full object-cover"
          loading={priority ? "eager" : "lazy"}
          decoding="async"
        />
      </div>
    </div>
  );
};

export default ImageCard;
