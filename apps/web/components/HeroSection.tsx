import { HeroSectionProps } from "@/types";

export default function HeroSection({ onGetStarted }: HeroSectionProps) {
  return (
    <section className="min-h-screen flex flex-col items-center justify-center text-center dark:bg-black bg-white px-4">
      <h1 className="text-4xl font-bold dark:text-purple-400 text-purple-700 mb-4">
        Welcome to Algo Buddy
      </h1>
      <p className="text-lg dark:text-gray-300 text-gray-700 max-w-xl mb-8">
        A collaborative platform where developers master DSA and tech stacks through peer learning, problem-solving, and real progress tracking.
      </p>
      <button
        onClick={onGetStarted}
        className="bg-purple-600 hover:bg-purple-700 text-white px-6 py-3 rounded-lg font-semibold"
      >
        Get Started
      </button>
    </section>
  );
}
