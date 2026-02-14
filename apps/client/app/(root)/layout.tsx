import Navbar from "../components/root/Navbar";

const layout = ({ children }: { children: React.ReactNode }) => {
  return (
    // root nav bar here
    <>
      <Navbar />
      {children}
    </>

    // root footer here
  );
};

export default layout;
