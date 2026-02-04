import { VisualizationData, VisualizationNode, TreeNode } from '../types';

/**
 * Converts flat visualization data to a hierarchical tree structure
 */
export function buildTreeHierarchy(data: VisualizationData): TreeNode | null {
  const nodeMap = new Map<string, TreeNode>();
  const childrenMap = new Map<string, string[]>();
  const spouseMap = new Map<string, string>();

  // Create node map
  data.nodes.forEach(node => {
    nodeMap.set(node.id, { ...node, children: [] });
  });

  // Build relationships from edges
  data.edges.forEach(edge => {
    if (edge.type === 'parent') {
      // edge.source is parent, edge.target is child
      const children = childrenMap.get(edge.source) || [];
      children.push(edge.target);
      childrenMap.set(edge.source, children);
    } else if (edge.type === 'spouse') {
      spouseMap.set(edge.source, edge.target);
      spouseMap.set(edge.target, edge.source);
    }
  });

  // Build tree starting from root
  const rootNode = nodeMap.get(data.root_id);
  if (!rootNode) return null;

  // Build descendant tree (going down from root)
  function buildDescendantTree(nodeId: string, visited: Set<string>): TreeNode | null {
    if (visited.has(nodeId)) return null;
    visited.add(nodeId);

    const node = nodeMap.get(nodeId);
    if (!node) return null;

    const treeNode: TreeNode = { ...node, children: [] };

    // Get children
    const childIds = childrenMap.get(nodeId) || [];
    childIds.forEach(childId => {
      const childNode = buildDescendantTree(childId, visited);
      if (childNode) {
        treeNode.children?.push(childNode);
      }
    });

    // Add spouse info
    const spouseId = spouseMap.get(nodeId);
    if (spouseId) {
      treeNode.spouse = nodeMap.get(spouseId);
    }

    return treeNode;
  }

  // For now, build descendant tree (root at top, descendants below)
  return buildDescendantTree(data.root_id, new Set());
}

/**
 * Calculate tree dimensions based on node count
 */
export function calculateTreeDimensions(nodeCount: number): { width: number; height: number } {
  const baseWidth = 1200;
  const baseHeight = 800;

  // Scale based on node count
  const scale = Math.max(1, Math.sqrt(nodeCount / 10));

  return {
    width: Math.min(baseWidth * scale, 3000),
    height: Math.min(baseHeight * scale, 2000)
  };
}

/**
 * Get color for a person based on gender
 */
export function getPersonColor(gender: 'M' | 'F', isAlive: boolean): string {
  if (!isAlive) {
    return gender === 'M' ? '#7a9cc6' : '#c6a7b8'; // Muted colors for deceased
  }
  return gender === 'M' ? '#4a90d9' : '#d94a7b'; // Vivid colors for living
}

/**
 * Format person display text
 */
export function formatPersonLabel(node: VisualizationNode): string {
  const years = node.death_year
    ? `${node.birth_year}-${node.death_year}`
    : `b. ${node.birth_year}`;
  return `${node.name}\n(${years})`;
}
