import { VisualizationData, VisualizationNode, TreeNode } from '../types';

  
                                                                     
  
export function buildTreeHierarchy(data: VisualizationData): TreeNode | null {
  const nodeMap = new Map<string, TreeNode>();
  const childrenMap = new Map<string, string[]>();
  const spouseMap = new Map<string, string>();

  
  data.nodes.forEach(node => {
    nodeMap.set(node.id, { ...node, children: [] });
  });

  
  data.edges.forEach(edge => {
    if (edge.type === 'parent') {
      
      const children = childrenMap.get(edge.source) || [];
      children.push(edge.target);
      childrenMap.set(edge.source, children);
    } else if (edge.type === 'spouse') {
      spouseMap.set(edge.source, edge.target);
      spouseMap.set(edge.target, edge.source);
    }
  });

  
  const rootNode = nodeMap.get(data.root_id);
  if (!rootNode) return null;

  
  function buildDescendantTree(nodeId: string, visited: Set<string>): TreeNode | null {
    if (visited.has(nodeId)) return null;
    visited.add(nodeId);

    const node = nodeMap.get(nodeId);
    if (!node) return null;

    const treeNode: TreeNode = { ...node, children: [] };

    
    const childIds = childrenMap.get(nodeId) || [];
    childIds.forEach(childId => {
      const childNode = buildDescendantTree(childId, visited);
      if (childNode) {
        treeNode.children?.push(childNode);
      }
    });

    
    const spouseId = spouseMap.get(nodeId);
    if (spouseId) {
      treeNode.spouse = nodeMap.get(spouseId);
    }

    return treeNode;
  }

  
  return buildDescendantTree(data.root_id, new Set());
}

  
                                                 
  
export function calculateTreeDimensions(nodeCount: number): { width: number; height: number } {
  const baseWidth = 1200;
  const baseHeight = 800;

  
  const scale = Math.max(1, Math.sqrt(nodeCount / 10));

  return {
    width: Math.min(baseWidth * scale, 3000),
    height: Math.min(baseHeight * scale, 2000)
  };
}

  
                                          
  
export function getPersonColor(gender: 'M' | 'F', isAlive: boolean): string {
  if (!isAlive) {
    return gender === 'M' ? '#7a9cc6' : '#c6a7b8'; 
  }
  return gender === 'M' ? '#4a90d9' : '#d94a7b'; 
}

  
                              
  
export function formatPersonLabel(node: VisualizationNode): string {
  const years = node.death_year
    ? `${node.birth_year}-${node.death_year}`
    : `b. ${node.birth_year}`;
  return `${node.name}\n(${years})`;
}
