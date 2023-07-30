const App = () => {
  return (
    <div className="bg-teal-900 h-screen w-full p-4">
      <div className="w-full flex justify-between items-center rounded bg-teal-800 p-4">
        <h1 className="text-3xl text-slate-300 font-bold">Golangify</h1>
        <div id="account_detail" className="text-bg-white bg-slate-900 p-2 rounded-md flex gap-2 items-center justify-center font-semibold">
          <div id="account_profile_picture" className="gap-4 w-8 h-8 rounded-full bg-slate-300"></div>
          <span className="text-slate-300">Name</span>
        </div>
      </div>
    </div>
  )
}

export default App
