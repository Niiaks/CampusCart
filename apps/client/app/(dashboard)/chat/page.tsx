"use client";

import { useState, useRef, useEffect } from "react";
import {
  Send,
  Paperclip,
  Smile,
  ChevronLeft,
  MoreVertical,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

// Dummy data for UI demo
const dummyConversations = [
  {
    id: "1",
    name: "Kwame Mensah",
    avatar: "/avatar1.png",
    lastMessage: "Thanks, I'll pick it up tomorrow!",
    time: "2m ago",
    unread: 2,
  },
  {
    id: "2",
    name: "Ama Boateng",
    avatar: "/avatar2.png",
    lastMessage: "Can you do 100 cedis?",
    time: "10m ago",
    unread: 0,
  },
  {
    id: "3",
    name: "John Doe",
    avatar: "/avatar3.png",
    lastMessage: "Sent you the payment!",
    time: "1h ago",
    unread: 0,
  },
];

const dummyMessages = [
  {
    id: 1,
    sender: "me",
    text: "Hi, is this still available?",
    time: "09:30",
  },
  {
    id: 2,
    sender: "them",
    text: "Yes, it is!",
    time: "09:31",
  },
  {
    id: 3,
    sender: "me",
    text: "Great! Can I pick it up tomorrow?",
    time: "09:32",
  },
  {
    id: 4,
    sender: "them",
    text: "Sure, what time works for you?",
    time: "09:33",
  },
  {
    id: 5,
    sender: "me",
    text: "Anytime after 2pm.",
    time: "09:34",
  },
  {
    id: 6,
    sender: "them",
    text: "Perfect, see you then!",
    time: "09:35",
  },
];

export default function ChatPage() {
  const [selected, setSelected] = useState(dummyConversations[0]);
  const [messages, setMessages] = useState(dummyMessages);
  const [input, setInput] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, selected]);

  const handleSend = () => {
    if (!input.trim()) return;
    setMessages((msgs) => [
      ...msgs,
      {
        id: msgs.length + 1,
        sender: "me",
        text: input,
        time: new Date().toLocaleTimeString([], {
          hour: "2-digit",
          minute: "2-digit",
        }),
      },
    ]);
    setInput("");
  };

  return (
    <div className="flex h-[calc(100vh-6rem)] w-full max-w-5xl mx-auto rounded-xl border border-border/60 bg-card shadow-lg overflow-hidden">
      {/* Sidebar */}
      <aside className="hidden w-80 shrink-0 border-r border-border/60 bg-background/80 md:flex flex-col">
        <div className="p-4 border-b border-border/60">
          <h2 className="text-lg font-bold text-foreground">Chats</h2>
        </div>
        <div className="flex-1 overflow-y-auto">
          {dummyConversations.map((conv) => (
            <button
              key={conv.id}
              onClick={() => setSelected(conv)}
              className={`flex w-full items-center gap-3 px-4 py-3 transition-colors hover:bg-muted/60 ${
                selected.id === conv.id ? "bg-brand/10" : ""
              }`}
            >
              <img
                src={conv.avatar}
                alt={conv.name}
                className="size-10 rounded-full object-cover border border-border"
              />
              <div className="min-w-0 flex-1 text-left">
                <p className="truncate font-medium text-foreground">
                  {conv.name}
                </p>
                <p className="truncate text-xs text-muted-foreground">
                  {conv.lastMessage}
                </p>
              </div>
              <div className="flex flex-col items-end gap-1">
                <span className="text-xs text-muted-foreground">
                  {conv.time}
                </span>
                {conv.unread > 0 && (
                  <span className="mt-1 inline-flex h-5 min-w-5 items-center justify-center rounded-full bg-brand text-xs font-bold text-brand-foreground px-1.5">
                    {conv.unread}
                  </span>
                )}
              </div>
            </button>
          ))}
        </div>
      </aside>

      {/* Chat area */}
      <section className="flex flex-1 flex-col">
        {/* Header */}
        <div className="flex items-center gap-3 border-b border-border/60 bg-background/80 px-4 py-3">
          <button className="md:hidden" onClick={() => {}}>
            <ChevronLeft className="size-5 text-muted-foreground" />
          </button>
          <img
            src={selected.avatar}
            alt={selected.name}
            className="size-10 rounded-full object-cover border border-border"
          />
          <div className="flex-1 min-w-0">
            <p className="truncate font-semibold text-foreground">
              {selected.name}
            </p>
            <span className="text-xs text-muted-foreground">Online</span>
          </div>
          <Button
            variant="ghost"
            size="icon-sm"
            className="text-muted-foreground"
          >
            <MoreVertical className="size-5" />
          </Button>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto bg-linear-to-b from-background/80 to-muted/60 px-4 py-6">
          <div className="flex flex-col gap-4">
            {messages.map((msg) => (
              <div
                key={msg.id}
                className={`flex ${msg.sender === "me" ? "justify-end" : "justify-start"}`}
              >
                <div
                  className={`max-w-xs rounded-2xl px-4 py-2 text-sm shadow-md ${
                    msg.sender === "me"
                      ? "bg-brand text-brand-foreground rounded-br-md"
                      : "bg-white text-foreground rounded-bl-md border border-border"
                  }`}
                >
                  {msg.text}
                  <span className="ml-2 align-bottom text-[10px] text-muted-foreground">
                    {msg.time}
                  </span>
                </div>
              </div>
            ))}
            <div ref={messagesEndRef} />
          </div>
        </div>

        {/* Input */}
        <form
          className="flex items-center gap-2 border-t border-border/60 bg-background/80 px-4 py-3"
          onSubmit={(e) => {
            e.preventDefault();
            handleSend();
          }}
        >
          <Button type="button" variant="ghost" size="icon-sm">
            <Smile className="size-5 text-muted-foreground" />
          </Button>
          <Button type="button" variant="ghost" size="icon-sm">
            <Paperclip className="size-5 text-muted-foreground" />
          </Button>
          <Input
            className="flex-1 rounded-full bg-muted px-4 text-sm"
            placeholder="Type a message..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            autoComplete="off"
          />
          <Button
            type="submit"
            size="icon-sm"
            className="bg-brand text-brand-foreground hover:bg-brand-hover"
            disabled={!input.trim()}
          >
            <Send className="size-5" />
          </Button>
        </form>
      </section>
    </div>
  );
}
