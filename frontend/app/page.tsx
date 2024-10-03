'use client';

import { useState } from "react";
import Header from "./components/Header";
import MainContent from "./components/MainContent";

export default function Home() {
  const [inputValue, setInputValue] = useState("");

  return (
    <div className="font-sans">
      {/* Header */}
      <Header />

      {/* Main content */}
      <MainContent />
    </div>
  );
}
