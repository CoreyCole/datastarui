package toast

// Toast expressions use the simple $toastID.open boolean model.
// No array/filter logic needed - each toast is an independent signal.
//
// Expression builders are in toast.templ:
//   - ShowToastExpr(id, duration) - exported for direct use in ClientActions
