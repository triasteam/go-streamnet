package tipselection

func UnIterableMap<HashId, Integer> calculate(Hash entryPoint) throws Exception {
	log.debug("Start calculating cw starting with tx hash {}", entryPoint);

	UnIterableMap<HashId, Integer> ret = new TransformingMap<>(HashPrefix::createPrefix, null);

	Set<Hash> visited = new HashSet<Hash> ();
	LinkedList<Hash> queue = new LinkedList<>();
	queue.add(entryPoint);
	Hash h;
	while (!queue.isEmpty()) {
		h = queue.pop();
		for (Hash e : tangle.getChild(h)) {
			if (tangle.contains(e) && !visited.contains(e)) {
				queue.add(e);
				visited.add(e);
			}
		}
		ret.put(h, (tangle.getScore(h).intValue()));
	}

	return ret;
}